package history

import (
	"golang.org/x/net/context"
	"time"

	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/reporter/view"
	statestream "github.com/fuserobotics/statestream"
)

type HistoryQueryBackend interface {
	// - copied from statestream.StorageBackend -
	// Retrieve the first snapshot before timestamp. Return nil for no data.
	GetSnapshotBefore(timestamp time.Time) (*statestream.StreamEntry, error)
	// Get the next entry after the timestamp. Return nil for no data.
	// Filter by the filter type, or don't filter if StreamEntryAny
	GetEntryAfter(timestamp time.Time, filterType statestream.StreamEntryType) (*statestream.StreamEntry, error)
	// Store a stream entry.
	SaveEntry(entry *statestream.StreamEntry) error
	// Amend an old entry
	AmendEntry(entry *statestream.StreamEntry, oldTimestamp time.Time) error
	// - additional -
	InitialSetComplete()
	// Provide a handle to the state stream
	SetSendStream(stream *statestream.Stream)
}

func HandleHistoryQuery(context context.Context, req *view.StateHistoryRequest, stream *statestream.Stream, backend HistoryQueryBackend) error {
	// Build a new stream for the send
	var streamConfig statestream.Config
	if req.StreamConfig != nil && req.StreamConfig.RecordRate != nil {
		streamConfig = *req.StreamConfig
	} else {
		streamConfig = stream.GetConfig()
	}

	sendStream, err := statestream.NewStream(backend, &streamConfig)
	if err != nil {
		return err
	}
	backend.SetSendStream(sendStream)

	var cursor *statestream.Cursor

	// first send the history
	earlyBound := util.NumberToTime(req.Query.BeginTime)
	lateBound := util.NumberToTime(req.Query.EndTime)
	if req.Query.EndTime == 0 {
		lateBound = time.Now()
	}

	cursor = stream.BuildCursor(statestream.ReadForwardCursor)
	ch := make(chan *statestream.StreamEntry)
	defer close(ch)

	go func() {
		for {
			select {
			case entry, ok := <-ch:
				if !ok {
					return
				}
				if entry == nil {
					backend.InitialSetComplete()
				} else {
					sendStream.WriteEntry(entry)
				}
			case <-context.Done():
				return
			}
		}
	}()

	if err := cursor.Init(earlyBound); err != nil {
		if err == statestream.NoDataError {
			// Check if we can get something after it.
			entry, err := stream.GetStorage().GetEntryAfter(earlyBound, statestream.StreamEntrySnapshot)
			if err != nil {
				return err
			}
			if entry == nil || entry.Timestamp.After(lateBound) {
				return statestream.NoDataError
			}
			earlyBound = entry.Timestamp
			if err := cursor.InitWithSnapshot(entry); err != nil {
				return err
			}
			// Send the initial snapshot entry
			select {
			case ch <- entry:
			case <-context.Done():
				return nil
			}
		} else {
			return err
		}
	}

	sub := cursor.SubscribeEntries(ch)

	cursor.SetTimestamp(lateBound)

	// close the history channel
	err = cursor.ComputeState()
	sub.Unsubscribe()

	if err != nil {
		return err
	}

	// if we need to tail, grab the write cursor.
	if !req.Query.Tail {
		return nil
	}

	cursor, err = stream.WriteCursor()
	if err != nil {
		return err
	}
	if !cursor.Ready() {
		if err := cursor.Init(time.Now()); err != nil {
			return err
		}
	}

	// Signal to the client the initial set is done
	ch <- nil

	sub = cursor.SubscribeEntries(ch)
	defer func() {
		sub.Unsubscribe()
	}()
	<-context.Done()
	return nil
}
