package reporter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/statestream"
	"github.com/golang/protobuf/proto"
)

var stateDetailsKey []byte = []byte("state")
var entryBucketName []byte = []byte("entries")

type State struct {
	hasEnsuredBuckets bool
	db                *bolt.DB
	stream            *stream.Stream
	dirtyChan         chan *State
	otherChan         chan *stream.StreamEntry

	Component  *Component
	Name       string
	BucketName []byte
	Data       dbproto.State
	OtherSub   stream.CursorEntrySubscription
}

func (s *State) Stream() *stream.Stream {
	if s.stream != nil {
		return s.stream
	}
	nstream, err := stream.NewStream(s, s.Data.StreamConfig)
	if err == nil {
		s.stream = nstream
	}
	return s.stream
}

func (s *State) getBucket(tx *bolt.Tx) *bolt.Bucket {
	cbkt := s.Component.getBucket(tx)
	sbkt := cbkt.Bucket(s.BucketName)
	return sbkt
}

func (s *State) getEntryBucket(tx *bolt.Tx) *bolt.Bucket {
	sbkt := s.getBucket(tx)
	return sbkt.Bucket(entryBucketName)
}

func (s *State) marshal() ([]byte, error) {
	return proto.Marshal(&s.Data)
}

func (s *State) LoadFromDb() error {
	return s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getBucket(tx)
		data := bkt.Get(stateDetailsKey)
		return proto.Unmarshal(data, &s.Data)
	})
}

// Slave this state to follow changes of another.
func (s *State) SubscribeToOther(other *State) error {
	cursor, err := other.Stream().WriteCursor()
	if err != nil {
		return err
	}

	if s.OtherSub != nil {
		s.OtherSub.Unsubscribe()
		close(s.otherChan)
		s.OtherSub = nil
	}

	// bufer this one to avoid write delays
	s.otherChan = make(chan *stream.StreamEntry, 50)
	s.OtherSub = cursor.SubscribeEntries(s.otherChan)
	go func() {
		for {
			select {
			case entry, ok := <-s.otherChan:
				if !ok {
					return
				}
				fmt.Printf("Writing locally from other stream: %v\n", entry)
				if err := s.WriteEntry(entry); err != nil {
					fmt.Printf("Error writing: %v\n", err)
				}
			}
		}
	}()
	return nil
}

func (s *State) Backfill(other *State) error {
	// Check if other is empty
	if len(other.Data.AllTimestamp) == 0 {
		return nil
	}

	// get stream
	str := s.Stream()

	// Check the latest point to backfill
	lastEntry := s.Data.LatestTimestamp()
	lastTargetEntry := other.Data.LatestTimestamp()

	for lastEntry < lastTargetEntry {
		ent, err := other.GetEntryAfter(util.NumberToTime(lastEntry), stream.StreamEntryAny)
		if err != nil {
			return err
		}

		if err := str.WriteEntry(ent); err != nil {
			return err
		}

		lastEntry = util.TimeToNumber(ent.Timestamp)
	}
	return nil
}

func (s *State) WriteState(timestamp time.Time, state stream.StateData) error {
	if err := s.Stream().WriteState(timestamp, state); err != nil {
		return err
	}
	if s.dirtyChan != nil {
		s.dirtyChan <- s
	}
	return nil
}

func (s *State) WriteEntry(entry *stream.StreamEntry) error {
	if err := s.Stream().WriteEntry(entry); err != nil {
		return err
	}
	if s.dirtyChan != nil {
		s.dirtyChan <- s
	}
	return nil
}

func (s *State) PurgeBefore(idx int) error {
	// first, commit the removal from the timestamps
	if idx < len(s.Data.AllTimestamp) {
		deletingBefore := s.Data.AllTimestamp[idx]
		// delete up to and including deletingIndex
		deletingIndex := -1
		for i, ts := range s.Data.SnapshotTimestamp {
			if ts >= deletingBefore {
				break
			}
			deletingIndex = i
		}
		if deletingIndex != -1 {
			s.Data.SnapshotTimestamp = s.Data.SnapshotTimestamp[:deletingIndex+1]
		}
	}
	tsToDelete := s.Data.AllTimestamp[:idx]
	s.Data.AllTimestamp = s.Data.AllTimestamp[idx:]
	if err := s.WriteToDb(); err != nil {
		return err
	}
	// now, actually delete the data
	return s.db.Update(func(tx *bolt.Tx) error {
		entryBkt := s.getEntryBucket(tx)
		// although we will pay attention to the error, don't let it stop us
		var rerr error
		for _, ts := range tsToDelete {
			if err := entryBkt.Delete([]byte(strconv.FormatInt(ts, 10))); err != nil {
				rerr = err
			}
		}
		return rerr
	})
}

func (s *State) writeToDbWithTransaction(tx *bolt.Tx) error {
	bkt := s.getBucket(tx)
	data, err := s.marshal()
	if err != nil {
		return err
	}
	if err := bkt.Put(stateDetailsKey, data); err != nil {
		return err
	}
	if s.hasEnsuredBuckets {
		return nil
	}
	if _, err := bkt.CreateBucketIfNotExists(entryBucketName); err != nil {
		return err
	}
	s.hasEnsuredBuckets = true
	return nil
}

func (s *State) WriteToDb() error {
	return s.db.Update(s.writeToDbWithTransaction)
}

func (s *State) PurgeFromDb() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		cbkt := s.Component.getBucket(tx)
		return cbkt.DeleteBucket(s.BucketName)
	})
}

func (s *State) Dispose() {
	if s.OtherSub != nil {
		s.OtherSub.Unsubscribe()
		close(s.otherChan)
		s.OtherSub = nil
	}
}
