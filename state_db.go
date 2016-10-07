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
	if snapshotCount == 0 {
		return nil, nil
	}
	idx := sort.Search(snapshotCount, func(i int) bool {
		return s.Data.SnapshotTimestamp[snapshotCount-i-1] < TimeToNumber(timestamp)
	})
	if idx < 0 || idx >= snapshotCount {
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
		Timestamp: NumberToTime(snapshotTime),
	}
	return entry, err
}

// Retrieve the first mutation after timestamp. Return nil for no data.
func (s *State) GetMutationAfter(timestamp time.Time) (*stream.StreamEntry, error) {
	mutationCount := len(s.Data.MutationTimestamp)
	if mutationCount == 0 {
		return nil, nil
	}
	idx := sort.Search(mutationCount, func(i int) bool {
		return s.Data.MutationTimestamp[i] > TimeToNumber(timestamp)
	})
	if idx < 0 || idx >= mutationCount {
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
		Timestamp: NumberToTime(mutationTime),
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
		timeNum := TimeToNumber(entry.Timestamp)
		switch entry.Type {
		case stream.StreamEntrySnapshot:
			bkt = s.getSnapshotBucket(tx)
			s.Data.SnapshotTimestamp = append(s.Data.SnapshotTimestamp, timeNum)
		case stream.StreamEntryMutation:
			bkt = s.getMutationBucket(tx)
			s.Data.MutationTimestamp = append(s.Data.MutationTimestamp, timeNum)
		}
		if err := bkt.Put([]byte(strconv.FormatInt(TimeToNumber(entry.Timestamp), 10)), data); err != nil {
			return err
		}
		return s.writeToDbWithTransaction(tx)
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
		return bkt.Put([]byte(strconv.FormatInt(TimeToNumber(oldTimestamp), 10)), data)
	})
}
