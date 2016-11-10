package reporter

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

var localBucketName []byte = []byte("local")

// A reporter instance
type Reporter struct {
	dbPath         string
	db             *bolt.DB
	hostIdentifier string
	enableRemotes  bool

	LocalTree  *ComponentTree
	RemoteList *RemoteList
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

func NewReporter(hostIdentifier, dbPath string, enableRemotes bool) (*Reporter, error) {
	res := &Reporter{dbPath: dbPath, hostIdentifier: hostIdentifier, enableRemotes: enableRemotes}
	if err := res.openDb(); err != nil {
		return nil, err
	}
	res.LocalTree = NewComponentTree(res.db, localBucketName)
	if err := res.LocalTree.loadComponentList(); err != nil {
		return nil, err
	}
	res.RemoteList = &RemoteList{
		db:         res.db,
		reporter:   res,
		BucketName: []byte("remote"),
		Remotes:    make(map[string]*Remote),
	}
	if err := res.RemoteList.LoadFromDb(); err != nil {
		return nil, err
	}
	if err := res.RemoteList.LoadAllRemotes(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Reporter) Close() {
	r.db.Close()
	r.LocalTree = nil
}
