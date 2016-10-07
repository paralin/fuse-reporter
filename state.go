package reporter

import (
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

	Component  *Component
	Name       string
	BucketName []byte
	Data       dbproto.State
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
