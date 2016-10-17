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

var entryTypeKey string = "$t"

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
	return s.getEntry(s.Data.SnapshotTimestamp[idx])
}

func (s *State) getEntry(timestamp int64) (*stream.StreamEntry, error) {
	var streamEntry interface{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getEntryBucket(tx)
		data := bkt.Get([]byte(strconv.FormatInt(timestamp, 10)))
		if len(data) == 0 {
			return errors.New("Data for timestamp not found.")
		}
		return json.Unmarshal(data, &streamEntry)
	})
	if err != nil {
		return nil, err
	}
	streamEntryData := streamEntry.(map[string]interface{})
	entryType := stream.StreamEntryType(streamEntryData[entryTypeKey].(float64))
	delete(streamEntryData, entryTypeKey)
	entry := &stream.StreamEntry{
		Type:      entryType,
		Data:      streamEntryData,
		Timestamp: NumberToTime(timestamp),
	}
	return entry, err
}

// Retrieve the first entry after timestamp. Return nil for no data.
func (s *State) GetEntryAfter(timestamp time.Time, filterType stream.StreamEntryType) (*stream.StreamEntry, error) {
	var tsArr []int64
	if filterType == stream.StreamEntrySnapshot {
		tsArr = s.Data.SnapshotTimestamp
	} else {
		tsArr = s.Data.AllTimestamp
	}
	entryCount := len(tsArr)
	if entryCount == 0 {
		return nil, nil
	}
	timeNum := TimeToNumber(timestamp)
	idx := sort.Search(entryCount, func(i int) bool {
		return tsArr[i] > timeNum
	})
	if idx < 0 || idx >= entryCount {
		return nil, nil
	}
	return s.getEntry(tsArr[idx])
}

func (s *State) writeEntryToDb(entry *stream.StreamEntry, isReplacement bool) error {
	jsonData := map[string]interface{}(entry.Data)
	jsonData[entryTypeKey] = int(entry.Type)
	data, err := json.Marshal(&jsonData)
	delete(jsonData, entryTypeKey)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := s.getEntryBucket(tx)
		timeNum := TimeToNumber(entry.Timestamp)
		if !isReplacement {
			s.Data.AllTimestamp = append(s.Data.AllTimestamp, timeNum)
			if entry.Type == stream.StreamEntrySnapshot {
				s.Data.SnapshotTimestamp = append(s.Data.SnapshotTimestamp, timeNum)
			}
		}
		if err := bkt.Put([]byte(strconv.FormatInt(TimeToNumber(entry.Timestamp), 10)), data); err != nil {
			return err
		}
		if isReplacement {
			return nil
		}
		return s.writeToDbWithTransaction(tx)
	})
}

func (s *State) pushEntryToRemotes(entry *stream.StreamEntry) {
	for _, rem := range s.RemoteStates {
		stream := rem.Stream()
		if stream == nil {
			continue
		}
		stream.WriteEntry(entry)
	}
}

// Store a stream entry.
func (s *State) SaveEntry(entry *stream.StreamEntry) error {
	s.pushEntryToRemotes(entry)
	return s.writeEntryToDb(entry, false)
}

// Amend an old entry
func (s *State) AmendEntry(entry *stream.StreamEntry, oldTimestamp time.Time) error {
	s.pushEntryToRemotes(entry)
	return s.writeEntryToDb(entry, true)
}
