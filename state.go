package reporter

import (
	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/fuserobotics/statestream"
	"github.com/golang/protobuf/proto"
)

var stateDetailsKey []byte = []byte("state")
var snapshotBucketName []byte = []byte("snapshots")
var mutationBucketName []byte = []byte("mutations")

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

func (s *State) getSnapshotBucket(tx *bolt.Tx) *bolt.Bucket {
	sbkt := s.getBucket(tx)
	return sbkt.Bucket(snapshotBucketName)
}

func (s *State) getMutationBucket(tx *bolt.Tx) *bolt.Bucket {
	sbkt := s.getBucket(tx)
	return sbkt.Bucket(mutationBucketName)
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

func (s *State) WriteToDb() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := s.getBucket(tx)
		data, err := s.marshal()
		if err != nil || s.hasEnsuredBuckets {
			return err
		}
		if err := bkt.Put(stateDetailsKey, data); err != nil {
			return err
		}
		if _, err := bkt.CreateBucketIfNotExists(snapshotBucketName); err != nil {
			return err
		}
		if _, err := bkt.CreateBucketIfNotExists(mutationBucketName); err != nil {
			return err
		}
		s.hasEnsuredBuckets = true
		return nil
	})
}
