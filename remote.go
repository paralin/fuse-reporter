package reporter

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/golang/protobuf/proto"
)

type Remote struct {
	db *bolt.DB

	RemoteList *RemoteList

	Data          dbproto.Remote
	DetailsKey    []byte
	ComponentTree *ComponentTree
}

func (c *Remote) LoadComponentTree() error {
	if c.Data.Id == "" {
		return errors.New("Remote ID cannot be empty.")
	}
	bktName := []byte(fmt.Sprintf("rem.%s", c.Data.Id))
	err := c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bktName)
		return err
	})
	if err != nil {
		return err
	}
	c.ComponentTree = NewComponentTree(c.db, bktName)
	if err := c.ComponentTree.loadComponentList(); err != nil {
		return err
	}
	return nil
}

func (c *Remote) LoadFromDb() error {
	return c.db.View(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		data := bkt.Get(c.DetailsKey)
		return proto.Unmarshal(data, &c.Data)
	})
}

func (c *Remote) WriteToDb() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		data, err := proto.Marshal(&c.Data)
		if err != nil {
			return err
		}
		return bkt.Put(c.DetailsKey, data)
	})
}

func (c *Remote) getBucket(tx *bolt.Tx) *bolt.Bucket {
	return c.RemoteList.getBucket(tx)
}
