package reporter

import (
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
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

	Component  *Component
	Name       string
	BucketName []byte
	Data       dbproto.State

	RemoteStates map[string]*State
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
		ent, err := other.GetEntryAfter(NumberToTime(lastEntry), stream.StreamEntryAny)
		if err != nil {
			return err
		}

		if err := str.WriteEntry(ent); err != nil {
			return err
		}

		lastEntry = TimeToNumber(ent.Timestamp)
	}
	return nil
}

func (s *State) WriteState(timestamp time.Time, state stream.StateData) error {
	if err := s.Stream().WriteState(timestamp, state); err != nil {
		return err
	}
	for _, rem := range s.RemoteStates {
		rem.WriteState(timestamp, state)
	}
	if s.dirtyChan != nil {
		s.dirtyChan <- s
	}
	return nil
}

func (s *State) PurgeBefore(idx int) error {
	// first, commit the removal from the timestamps
	tsToDelete := s.Data.AllTimestamp[:idx]
	s.Data.AllTimestamp = s.Data.AllTimestamp[idx:]
	if err := s.WriteToDb(); err != nil {
		return err
	}
	// now, actually delete the data
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := s.getBucket(tx)
		// although we will pay attention to the error, don't let it stop us
		var rerr error
		for _, ts := range tsToDelete {
			if err := bkt.Delete([]byte(strconv.FormatInt(ts, 10))); err != nil {
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
