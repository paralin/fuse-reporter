package reporter

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/golang/protobuf/proto"
)

var ComponentNotExistError error = errors.New("Component does not exist.")

// Top-level list of components
var componentListKey []byte = []byte("components")
var componentDetailsKey []byte = []byte("component")

type ComponentList struct {
	Data dbproto.ComponentList
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
		bkt := tx.Bucket(globalBucketName)
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

// A instantiated component
type Component struct {
	Name          string
	BucketName    []byte
	ComponentList *ComponentList
	Data          dbproto.Component
}

func (c *Component) getBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(c.BucketName)
}

func (c *Component) LoadFromDb(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		data := bkt.Get(componentDetailsKey)
		return proto.Unmarshal(data, &c.Data)
	})
}

func LoadComponent(db *bolt.DB, componentList *ComponentList, name string) (*Component, error) {
	if !componentList.ContainsComponent(name) {
		return nil, ComponentNotExistError
	}
	res := &Component{Name: name, ComponentList: componentList, BucketName: []byte(fmt.Sprintf("cmp.%s", name))}
	return res, res.LoadFromDb(db)
}
