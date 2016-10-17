package reporter

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/golang/protobuf/proto"
)

var remoteListKey []byte = []byte("remotes")

type RemoteList struct {
	db         *bolt.DB
	BucketName []byte
	Data       dbproto.RemoteList

	Remotes map[string]*Remote
}

func (c *RemoteList) getBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(c.BucketName)
}

func (c *RemoteList) WriteToDb() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(c.BucketName)
		data, err := proto.Marshal(&c.Data)
		if err != nil {
			return err
		}
		return bkt.Put(remoteListKey, data)
	})
}

func (c *RemoteList) ContainsRemote(remoteId string) bool {
	for _, rid := range c.Data.RemoteId {
		if rid == remoteId {
			return true
		}
	}
	return false
}

func (c *RemoteList) LoadFromDb() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(c.BucketName)
		if err != nil {
			return err
		}
		data := bkt.Get(remoteListKey)
		c.Data.Reset()
		if data != nil {
			if err := proto.Unmarshal(data, &c.Data); err != nil {
				return err
			}
		}
		dta, err := proto.Marshal(&c.Data)
		if err != nil {
			return err
		}
		return bkt.Put(remoteListKey, dta)
	})
}

func (c *RemoteList) LoadAllRemotes() error {
	for _, id := range c.Data.RemoteId {
		rem := &Remote{
			db:         c.db,
			RemoteList: c,
			DetailsKey: []byte(fmt.Sprintf("r.%s", id)),
		}
		rem.Manager = &RemoteManager{r: rem}
		if err := rem.Init(); err != nil {
			return err
		}
		c.Remotes[id] = rem
	}
	return nil
}

func (r *RemoteList) CreateOrUpdateRemote(remoteId, endpoint string) (*Remote, error) {
	exist, ok := r.Remotes[remoteId]
	if ok {
		if exist.Data.Config.Endpoint == endpoint {
			return exist, nil
		}
	}
	nrem := &Remote{
		db:         r.db,
		RemoteList: r,
		DetailsKey: []byte(fmt.Sprintf("r.%s", remoteId)),
		Data: dbproto.Remote{
			Id: remoteId,
			Config: &dbproto.RemoteConfig{
				Endpoint: endpoint,
			},
		},
	}
	if err := nrem.WriteToDb(); err != nil {
		return nil, err
	}
	r.Remotes[remoteId] = nrem
	r.Data.RemoteId = append(r.Data.RemoteId, remoteId)
	if err := r.WriteToDb(); err != nil {
		return nil, err
	}
	if err := nrem.LoadComponentTree(); err != nil {
		return nil, err
	}
	if err := nrem.ComponentTree.ComponentList.WriteToDb(r.db); err != nil {
		return nil, err
	}
	return nrem, nil
}
