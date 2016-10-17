package reporter

import (
	"github.com/boltdb/bolt"
)

type ComponentTree struct {
	db *bolt.DB

	BucketName []byte

	Components    map[string]*Component
	ComponentList *ComponentList
}

func NewComponentTree(db *bolt.DB, bucketName []byte) *ComponentTree {
	return &ComponentTree{
		db:         db,
		BucketName: bucketName,
		Components: make(map[string]*Component),
	}
}

func (c *ComponentTree) getBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(c.BucketName)
}

func (r *ComponentTree) GetComponent(componentName string) (*Component, error) {
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
	r.Components[componentName] = component
	return component, nil
}

func (r *ComponentTree) CreateComponentIfNotExists(componentName string) (*Component, error) {
	if component, err := r.GetComponent(componentName); err == nil {
		return component, err
	}

	component, err := CreateComponent(r.db, r.ComponentList, componentName)
	if err != nil {
		return nil, err
	}
	r.Components[componentName] = component
	return component, nil
}

func (r *ComponentTree) loadComponentList() error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(r.BucketName)
		return nil
	})
	if err != nil {
		return err
	}
	r.ComponentList = &ComponentList{BucketName: r.BucketName}
	return r.ComponentList.LoadFromDb(r.db)
}

func (r *ComponentTree) getAllComponents() []*Component {
	res := make([]*Component, len(r.ComponentList.Data.ComponentName))
	for i, comp := range r.ComponentList.Data.ComponentName {
		c, _ := r.GetComponent(comp)
		res[i] = c
	}
	return res
}
