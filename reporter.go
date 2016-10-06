package reporter

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

var globalBucketName []byte = []byte("global")

// A reporter instance
type Reporter struct {
	dbPath string
	db     *bolt.DB

	Components    map[string]*Component
	ComponentList *ComponentList
}

func (r *Reporter) openDb() error {
	glog.Infof("Attempting to open DB at %s...", r.dbPath)
	db, err := bolt.Open(r.dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	glog.Infof("DB opened successfully.")
	r.db = db
	return nil
}

func (r *Reporter) GetComponent(componentName string) (*Component, error) {
	component, ok := r.Components[componentName]
	if ok {
		return component, nil
	}

	// Attempt to load it from DB
	component, err := LoadComponent(r.db, r.ComponentList, componentName)
	if err != nil {
		return nil, err
	}
	if component == nil {
		return nil, ComponentNotExistError
	}
	return component, nil
}

func (r *Reporter) loadComponentList() error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(globalBucketName)
		return nil
	})
	if err != nil {
		return err
	}
	r.ComponentList = &ComponentList{}
	return r.ComponentList.LoadFromDb(r.db)
}

func NewReporter(dbPath string) (*Reporter, error) {
	res := &Reporter{dbPath: dbPath, Components: make(map[string]*Component)}
	if err := res.openDb(); err != nil {
		return nil, err
	}
	if err := res.loadComponentList(); err != nil {
		return nil, err
	}
	return nil, nil
}
