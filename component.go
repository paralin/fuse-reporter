package reporter

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fuserobotics/reporter/dbproto"
	"github.com/fuserobotics/statestream"
	"github.com/golang/protobuf/proto"
)

var ComponentNotExistError error = errors.New("Component or state does not exist.")

var componentDetailsKey []byte = []byte("component")

// A instantiated component
type Component struct {
	db        *bolt.DB
	dirtyChan chan *State

	Name          string
	BucketName    []byte
	ComponentList *ComponentList
	Data          dbproto.Component

	States map[string]*State
}

func (c *Component) getBucket(tx *bolt.Tx) *bolt.Bucket {
	pbkt := c.ComponentList.getBucket(tx)
	return pbkt.Bucket(c.BucketName)
}

func (c *Component) getAllStates() []*State {
	res := make([]*State, len(c.Data.StateName))
	for i, sn := range c.Data.StateName {
		st, _ := c.GetState(sn)
		res[i] = st
	}
	return res
}

func (c *Component) LoadFromDb(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		data := bkt.Get(componentDetailsKey)
		return proto.Unmarshal(data, &c.Data)
	})
}

func (c *Component) hasState(stateName string) bool {
	for _, sn := range c.Data.StateName {
		if sn == stateName {
			return true
		}
	}
	return false
}

func (c *Component) CreateStateIfNotExists(stateName string, config *stream.Config, initialEntry *stream.StreamEntry) (*State, error) {
	state, err := c.GetState(stateName)
	if err == nil {
		return state, nil
	}
	state = &State{
		db:         c.db,
		Component:  c,
		Name:       stateName,
		BucketName: []byte(fmt.Sprintf("st.%s", stateName)),
		dirtyChan:  c.dirtyChan,
	}
	state.Data.StreamConfig = config
	c.db.Update(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		bkt.CreateBucketIfNotExists(state.BucketName)
		return nil
	})
	c.States[stateName] = state
	c.Data.StateName = append(c.Data.StateName, stateName)
	c.WriteToDb()
	if err := state.WriteToDb(); err != nil {
		return nil, err
	}
	if initialEntry != nil {
		if err := state.WriteEntry(initialEntry); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func (c *Component) ReleaseAllStates() {
	for id, state := range c.States {
		delete(c.States, id)
		state.Dispose()
	}
}

func (c *Component) DeleteState(stateName string) error {
	state, err := c.GetState(stateName)
	if err == ComponentNotExistError {
		return nil
	}
	delete(c.States, stateName)
	state.Dispose()
	// delete from name list
	for idx, nm := range c.Data.StateName {
		if nm == stateName {
			c.Data.StateName = append(c.Data.StateName[:idx], c.Data.StateName[idx+1:]...)
			if err := c.WriteToDb(); err != nil {
				return err
			}
		}
	}
	// purge from db
	return state.PurgeFromDb()
}

func (c *Component) GetState(stateName string) (*State, error) {
	state, ok := c.States[stateName]
	if ok {
		return state, nil
	}

	if !c.hasState(stateName) {
		return nil, ComponentNotExistError
	}

	// Attempt to load it from DB
	state = &State{
		Component:  c,
		Name:       stateName,
		BucketName: []byte(fmt.Sprintf("st.%s", stateName)),
		db:         c.db,
		dirtyChan:  c.dirtyChan,
	}
	if err := state.LoadFromDb(); err != nil {
		return nil, err
	}
	c.States[stateName] = state
	return state, nil
}

func (c *Component) marshal() ([]byte, error) {
	return proto.Marshal(&c.Data)
}

func (c *Component) WriteToDb() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		bkt := c.getBucket(tx)
		dta, err := c.marshal()
		if err != nil {
			return err
		}
		return bkt.Put(componentDetailsKey, dta)
	})
}

func (c *Component) PurgeFromDb() error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(c.BucketName)
	})
}

func newComponent(db *bolt.DB, componentList *ComponentList, name string) *Component {
	return &Component{
		Name:          name,
		ComponentList: componentList,
		BucketName:    []byte(fmt.Sprintf("cmp.%s", name)),
		States:        make(map[string]*State),
		db:            db,
	}
}

func LoadComponent(db *bolt.DB, componentList *ComponentList, name string) (*Component, error) {
	if !componentList.ContainsComponent(name) {
		return nil, ComponentNotExistError
	}
	res := newComponent(db, componentList, name)
	return res, res.LoadFromDb(db)
}

func CreateComponent(db *bolt.DB, componentList *ComponentList, name string) (*Component, error) {
	res := newComponent(db, componentList, name)
	err := db.Update(func(tx *bolt.Tx) error {
		pbkt := componentList.getBucket(tx)
		bkt, err := pbkt.CreateBucketIfNotExists(res.BucketName)
		if err != nil {
			return err
		}
		dta, err := res.marshal()
		if err != nil {
			return err
		}
		return bkt.Put(componentDetailsKey, dta)
	})
	if err != nil {
		return res, err
	}
	componentList.Data.ComponentName = append(componentList.Data.ComponentName, name)
	return res, componentList.WriteToDb(db)
}
