package reporter

import (
	"errors"
	"time"

	"github.com/fuserobotics/reporter/remote"
	"github.com/golang/glog"
	"golang.org/x/net/context"
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
	glog.Infof("Starting remote manager for %s...", m.r.Data.Config.Endpoint)
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

func (m *RemoteManager) baseContext() *remote.RequestContext {
	return &remote.RequestContext{
		HostIdentifier: m.r.RemoteList.reporter.hostIdentifier,
	}
}

func (m *RemoteManager) queryConfig() error {
	resp, err := m.client.GetRemoteConfig(context.Background(), &remote.GetRemoteConfigRequest{
		Context: m.baseContext(),
	})
	if err != nil {
		return err
	}
	conf := resp.Config
	if conf == nil {
		return errors.New("GetRemoteConfig() returned nil config.")
	}
	m.r.Data.StreamConfig = conf
	return m.r.WriteToDb()
}

func (m *RemoteManager) update() bool {
	if m.client == nil {
		glog.Infof("Attempting to connect to %s...", m.r.Data.Config.Endpoint)
		if err := m.connect(); err != nil {
			glog.Warningf("Failed to connect to remote, %v.", err)
			return !m.shouldStopTimeout(time.Duration(5) * time.Second)
		}
	}
	if m.r.Data.StreamConfig == nil {
		glog.Infof("Querying %s for config...", m.r.Data.Config.Endpoint)
		if err := m.queryConfig(); err != nil {
			glog.Warningf("Failed to get config, %v.", err)
			m.conn.Close()
			m.conn = nil
			m.client = nil
			return !m.shouldStopTimeout(time.Duration(5) * time.Second)
		}
		glog.Infof("Got from %s config with crc %d.", m.r.Data.Config.Endpoint, m.r.Data.StreamConfig.Crc32)
	}
	return !m.shouldStopTimeout(time.Duration(5) * time.Second)
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
