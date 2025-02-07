package reporter

import (
	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/golang/protobuf/proto"
)

// Top-level list of components
var componentListKey []byte = []byte("components")

type ComponentList struct {
	BucketName []byte
	Data       dbproto.ComponentList
}

func (c *ComponentList) getBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(c.BucketName)
}

func (c *ComponentList) WriteToDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(c.BucketName)
		data, err := proto.Marshal(&c.Data)
		if err != nil {
			return err
		}
		return bkt.Put(componentListKey, data)
	})
}

func (c *ComponentList) ContainsComponent(name string) bool {
	for _, nm := range c.Data.ComponentName {
		if nm == name {
			return true
		}
	}
	return false
}

func (c *ComponentList) LoadFromDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(c.BucketName)
		data := bkt.Get(componentListKey)
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
		return bkt.Put(componentListKey, dta)
	})
}

func (c *ComponentList) remove(db *bolt.DB, name string) error {
	for idx, nm := range c.Data.ComponentName {
		if nm == name {
			c.Data.ComponentName = append(c.Data.ComponentName[:idx], c.Data.ComponentName[idx+1:]...)
			return c.WriteToDb(db)
		}
	}
	return nil
}
