package reporter

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/fuserobotics/statestream"
)

// Retrieve the first snapshot before timestamp. Return nil for no data.
func (s *State) GetSnapshotBefore(timestamp time.Time) (*stream.StreamEntry, error) {
	snapshotCount := len(s.Data.SnapshotTimestamp)
	idx := sort.Search(snapshotCount, func(i int) bool {
		return s.Data.SnapshotTimestamp[snapshotCount-i-1] < timestamp.UnixNano()
	})
	if idx < 0 {
		return nil, nil
	}
	snapshotTime := s.Data.SnapshotTimestamp[idx]
	var streamEntry interface{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getSnapshotBucket(tx)
		data := bkt.Get([]byte(strconv.FormatInt(snapshotTime, 10)))
		if len(data) == 0 {
			return errors.New("Data for timestamp not found.")
		}
		return json.Unmarshal(data, &streamEntry)
	})
	entry := &stream.StreamEntry{
		Type:      stream.StreamEntrySnapshot,
		Data:      streamEntry.(map[string]interface{}),
		Timestamp: time.Unix(0, snapshotTime),
	}
	return entry, err
}

// Retrieve the first mutation after timestamp. Return nil for no data.
func (s *State) GetMutationAfter(timestamp time.Time) (*stream.StreamEntry, error) {
	mutationCount := len(s.Data.MutationTimestamp)
	idx := sort.Search(mutationCount, func(i int) bool {
		return s.Data.MutationTimestamp[i] > timestamp.UnixNano()
	})
	if idx < 0 {
		return nil, nil
	}
	mutationTime := s.Data.MutationTimestamp[idx]
	var streamEntry interface{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getMutationBucket(tx)
		data := bkt.Get([]byte(strconv.FormatInt(mutationTime, 10)))
		if len(data) == 0 {
			return errors.New("Data for timestamp not found.")
		}
		return json.Unmarshal(data, &streamEntry)
	})
	entry := &stream.StreamEntry{
		Type:      stream.StreamEntryMutation,
		Data:      streamEntry.(map[string]interface{}),
		Timestamp: time.Unix(0, mutationTime),
	}
	return entry, err
}

// Store a stream entry.
func (s *State) SaveEntry(entry *stream.StreamEntry) error {
	jsonData := map[string]interface{}(entry.Data)
	data, err := json.Marshal(&jsonData)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		var bkt *bolt.Bucket
		switch entry.Type {
		case stream.StreamEntrySnapshot:
			bkt = s.getSnapshotBucket(tx)
		case stream.StreamEntryMutation:
			bkt = s.getMutationBucket(tx)
		}
		return bkt.Put([]byte(strconv.FormatInt(entry.Timestamp.UnixNano(), 10)), data)
	})
}

// Amend an old entry
func (s *State) AmendEntry(entry *stream.StreamEntry, oldTimestamp time.Time) error {
	jsonData := map[string]interface{}(entry.Data)
	data, err := json.Marshal(&jsonData)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		var bkt *bolt.Bucket
		switch entry.Type {
		case stream.StreamEntrySnapshot:
			bkt = s.getSnapshotBucket(tx)
		case stream.StreamEntryMutation:
			bkt = s.getMutationBucket(tx)
		}
		return bkt.Put([]byte(strconv.FormatInt(oldTimestamp.UnixNano(), 10)), data)
	})
}
