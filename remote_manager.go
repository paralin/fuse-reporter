package reporter

import (
	"time"

	"github.com/fuserobotics/reporter/remote"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var connectionTimeout time.Duration = time.Duration(5) * time.Second

type RemoteManager struct {
	r *Remote

	conn   *grpc.ClientConn
	client remote.ReporterRemoteServiceClient

	stopChan chan bool
}

func (m *RemoteManager) Start() {
	m.stopChan = make(chan bool, 1)
	go m.updateThread()
}

func (m *RemoteManager) connect() error {
	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
		m.client = nil
	}

	conn, err := grpc.Dial(m.r.Data.Config.Endpoint, grpc.WithTimeout(connectionTimeout), grpc.WithInsecure())
	if err != nil {
		return err
	}
	m.conn = conn

	m.client = remote.NewReporterRemoteServiceClient(conn)
	return nil
}

func (m *RemoteManager) update() bool {
	if m.client == nil {
		glog.Infof("Attempting to connect to %s...", m.r.Data.Config.Endpoint)
		if err := m.connect(); err != nil {
			glog.Warningf("Failed to connect to remote, %v.", err)
			return !m.shouldStopTimeout(time.Duration(5) * time.Second)
		}
	}
	return true
}

func (m *RemoteManager) cleanup() {
	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}
}

func (m *RemoteManager) updateThread() {
	defer m.cleanup()
	for {
		if m.shouldStop() {
			break
		}

		if !m.update() {
			break
		}
	}
}

func (m *RemoteManager) shouldStop() bool {
	select {
	case <-m.stopChan:
		return true
	default:
		return false
	}
}

func (m *RemoteManager) shouldStopTimeout(dur time.Duration) bool {
	timer := time.NewTimer(dur)
	defer timer.Stop()
	select {
	case <-m.stopChan:
		return true
	case <-timer.C:
		return false
	}
}

func (m *RemoteManager) Destroy() {
	m.stopChan <- true
}
