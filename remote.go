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

	Data               dbproto.Remote
	DetailsKey         []byte
	ComponentTree      *ComponentTree
	LocalComponentTree *ComponentTree

	Manager *RemoteManager
}

// Load everything, start background manager
func (rem *Remote) Init() error {
	if err := rem.LoadFromDb(); err != nil {
		return err
	}
	if err := rem.LoadComponentTree(); err != nil {
		return err
	}
	rem.Manager.Start()
	return nil
}

// Shut down the remote manager, etc
func (rem *Remote) Destroy() {
	rem.Manager.Destroy()
}

func (c *Remote) initAllComponents() []error {
	components := c.ComponentTree.getAllComponents()
	rerr := []error{}
	for _, comp := range components {
		if comp == nil {
			continue
		}

		lc, err := c.LocalComponentTree.GetComponent(comp.Name)
		if err != nil {
			rerr = append(rerr, err)
		}
		if lc == nil {
			continue
		}

		states := comp.getAllStates()
		for _, state := range states {
			if state == nil {
				continue
			}

			ls, err := lc.GetState(state.Name)
			if err != nil {
				rerr = append(rerr, err)
			}
			if err := state.Backfill(ls); err != nil {
				rerr = append(rerr, err)
			}
			state.RemoteStates = append(state.RemoteStates, state)
		}
	}

	return rerr
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
	if err := c.initAllComponents(); len(err) != 0 {
		return err[0]
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
