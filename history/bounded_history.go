package history

import (
	"encoding/json"
	"time"

	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/reporter/view"
	statestream "github.com/fuserobotics/statestream"
)

func buildStateEntryResponse(entry *statestream.StreamEntry) *view.StateEntry {
	if entry == nil {
		return nil
	}
	data, _ := json.Marshal(entry.Data)
	return &view.StateEntry{
		JsonState: string(data),
		Timestamp: util.TimeToNumber(entry.Timestamp),
		Type:      int32(entry.Type),
	}
}

func HandleBoundedHistoryQuery(req *view.BoundedStateHistoryRequest, srvstream view.ReporterService_GetBoundedStateHistoryServer, stream *statestream.Stream) error {
	// First, get the snapshot before the bound.
	midTimestamp := util.NumberToTime(req.MidTimestamp)
	storage := stream.GetStorage()
	startEntry, err := storage.GetSnapshotBefore(midTimestamp)
	if err != nil {
		return err
	}

	// send off the entry
	if err := srvstream.Send(&view.BoundedStateHistoryResponse{
		State:  buildStateEntryResponse(startEntry),
		Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_START_BOUND,
	}); err != nil {
		return err
	}

	if startEntry == nil {
		return nil
	}

	// Grab end snapshot
	endEntry, err := storage.GetEntryAfter(startEntry.Timestamp, statestream.StreamEntrySnapshot)
	if err != nil {
		return err
	}

	if err := srvstream.Send(&view.BoundedStateHistoryResponse{
		State:  buildStateEntryResponse(endEntry),
		Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_END_BOUND,
	}); err != nil {
		return err
	}

	// Now, push the data we already have.
	doneChan := srvstream.Context().Done()
	lastEntryTimestamp := startEntry.Timestamp
	for {
		nextEntry, err := storage.GetEntryAfter(lastEntryTimestamp, statestream.StreamEntryAny)
		if err != nil {
			return err
		}
		if nextEntry == nil || (endEntry != nil && endEntry.Timestamp.Before(nextEntry.Timestamp)) {
			break
		}
		if nextEntry.Type == statestream.StreamEntrySnapshot {
			break
		}

		// Check if the client closed the connection
		select {
		case <-doneChan:
			return nil
		default:
		}

		// Send out the entry
		if err := srvstream.Send(&view.BoundedStateHistoryResponse{
			State:  buildStateEntryResponse(nextEntry),
			Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_INITIAL_SET,
		}); err != nil {
			return err
		}
		lastEntryTimestamp = nextEntry.Timestamp
	}

	// Do we need to tail?
	if endEntry != nil {
		return nil
	}

	// Signal to the client the initial set is done
	if err := srvstream.Send(&view.BoundedStateHistoryResponse{
		State:  nil,
		Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_TAIL,
	}); err != nil {
		return err
	}

	// if we need to tail, grab the write cursor.
	cursor, err := stream.WriteCursor()
	if err != nil {
		return err
	}
	if !cursor.Ready() {
		if err := cursor.Init(time.Now()); err != nil {
			return err
		}
	}

	ch := make(chan *statestream.StreamEntry)
	defer close(ch)

	sub := cursor.SubscribeEntries(ch)
	defer func() {
		sub.Unsubscribe()
	}()

	for {
		select {
		case <-doneChan:
			return nil
		case entry, ok := <-ch:
			if !ok {
				return nil
			}
			// Conclude when we hit a snapshot
			if entry.Type == statestream.StreamEntrySnapshot {
				return srvstream.Send(&view.BoundedStateHistoryResponse{
					State:  buildStateEntryResponse(entry),
					Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_END_BOUND,
				})
			}
			if err := srvstream.Send(&view.BoundedStateHistoryResponse{
				State:  buildStateEntryResponse(entry),
				Status: view.BoundedStateHistoryResponse_BOUNDED_HISTORY_TAIL,
			}); err != nil {
				return err
			}
		}
	}
}
