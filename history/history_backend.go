package history

import (
	"encoding/json"
	"time"

	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/reporter/view"
	statestream "github.com/fuserobotics/statestream"
)

// Pending entry to send
type stateHistoryEntry struct {
	entry      *statestream.StreamEntry
	initialSet bool
}

func (s *stateHistoryEntry) BuildResponse() (*view.StateHistoryResponse, error) {
	status := view.StateHistoryResponse_HISTORY_TAIL
	if s.initialSet {
		status = view.StateHistoryResponse_HISTORY_INITIAL_SET
	}

	var state *view.StateEntry
	if s.entry != nil {
		jsonData, err := json.Marshal(s.entry.Data)
		if err != nil {
			return nil, err
		}
		state = &view.StateEntry{
			JsonState: string(jsonData),
			Timestamp: util.TimeToNumber(s.entry.Timestamp),
			Type:      int32(s.entry.Type),
		}
	}
	return &view.StateHistoryResponse{Status: status, State: state}, nil
}

type stateHistoryBackend struct {
	// Buffer of entries to send out
	sendChan chan *stateHistoryEntry
	// Pending entry, might be amended. Not used during tailing stage.
	pendingEntry *statestream.StreamEntry
	// Currently in the initial set?
	initialSet bool
	// Reference to the stream
	sendStream *statestream.Stream
}

type stateHistoryRequestBackend struct {
	stateHistoryBackend
	srvstream view.ReporterService_GetStateHistoryServer
	req       *view.StateHistoryRequest
	stream    *statestream.Stream
}

func HandleStateHistoryRequest(req *view.StateHistoryRequest, srvstream view.ReporterService_GetStateHistoryServer, stream *statestream.Stream) error {
	back := &stateHistoryRequestBackend{
		stateHistoryBackend: stateHistoryBackend{
			sendChan:   make(chan *stateHistoryEntry, 10),
			initialSet: true,
		},
		srvstream: srvstream,
		req:       req,
		stream:    stream,
	}
	return back.Handle()
}

func (back *stateHistoryBackend) SetSendStream(stream *statestream.Stream) {
	back.sendStream = stream
}

func (back *stateHistoryRequestBackend) Handle() error {
	sendErr := make(chan error, 1)
	go func() (sendError error) {
		defer func() {
			sendErr <- sendError
		}()
		for {
			select {
			case <-back.srvstream.Context().Done():
				return
			case ent, ok := <-back.stateHistoryBackend.sendChan:
				if !ok {
					return nil
				}
				resp, err := ent.BuildResponse()
				if err != nil {
					return err
				}
				if err := back.srvstream.Send(resp); err != nil {
					return err
				}
			}
		}
	}()

	histErr := make(chan error, 1)
	go func() {
		histErr <- HandleHistoryQuery(back.srvstream.Context(), back.req, back.stream, &back.stateHistoryBackend)
	}()

	sendComplete := false
	for !sendComplete {
		select {
		// if we get an error sending messages
		case err := <-sendErr:
			if err != nil {
				return err
			}
			sendComplete = true
		// if the client terminates the connection
		case <-back.srvstream.Context().Done():
			return nil
		// if we get an error fetching the history
		case err := <-histErr:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *stateHistoryBackend) flushPendingEntry() {
	if b.pendingEntry == nil {
		return
	}
	pe := b.pendingEntry
	b.pendingEntry = nil
	b.sendChan <- &stateHistoryEntry{entry: pe, initialSet: b.initialSet}
}

func (b *stateHistoryBackend) GetSnapshotBefore(timestamp time.Time) (*statestream.StreamEntry, error) {
	return nil, nil
}

// Get the next entry after the timestamp. Return nil for no data.
// Filter by the filter type, or don't filter if StreamEntryAny
func (b *stateHistoryBackend) GetEntryAfter(timestamp time.Time, filterType statestream.StreamEntryType) (*statestream.StreamEntry, error) {
	return nil, nil
}

// Store a stream entry.
func (b *stateHistoryBackend) SaveEntry(entry *statestream.StreamEntry) error {
	b.flushPendingEntry()
	b.pendingEntry = entry
	if !b.initialSet {
		b.flushPendingEntry()
	}

	return nil
}

// Amend an old entry
func (b *stateHistoryBackend) AmendEntry(entry *statestream.StreamEntry, oldTimestamp time.Time) error {
	if b.initialSet {
		b.pendingEntry = entry
	}

	return nil
}

// - additional -
func (b *stateHistoryBackend) InitialSetComplete() {
	b.flushPendingEntry()
	b.initialSet = false
	// inform the client the set is complete
	b.sendChan <- &stateHistoryEntry{entry: nil, initialSet: false}
	// disable amends so we send mutations instantly
	b.sendStream.DisableAmends()
}
