package reporter

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fuserobotics/reporter/remote"
	"github.com/fuserobotics/statestream"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var connectionTimeout time.Duration = time.Duration(5) * time.Second

type RemoteManager struct {
	r *Remote

	conn   *grpc.ClientConn
	client remote.ReporterRemoteServiceClient

	stopChan    chan bool
	pendingPush map[*State]int
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

func (m *RemoteManager) flushDirtyChan(commit bool) {
	for {
		select {
		case st := <-m.r.StateDirtyChan:
			if commit {
				if _, ok := m.pendingPush[st]; ok {
					m.pendingPush[st]++
				} else {
					m.pendingPush[st] = 1
				}
			}
		default:
			return
		}
	}
}

func (m *RemoteManager) buildPushSlice() {
	res := map[*State]int{}
	for _, cmp := range m.r.ComponentTree.getAllComponents() {
		for _, st := range cmp.getAllStates() {
			if st.Data.LastTimestamp != 0 && st.Data.RemoteTimestamp < st.Data.LastTimestamp {
				res[st] = 1
			}
		}
	}
	m.pendingPush = res
	m.flushDirtyChan(false)
}

// Synchronize config with local component tree
func (m *RemoteManager) syncRemoteConfig() error {
	ct := m.r.ComponentTree
	lct := m.r.RemoteList.reporter.LocalTree
	m.pendingPush = nil

	componentStates := make(map[string]map[string]bool)
	regCmpState := func(component, state string) {
		cs, ok := componentStates[component]
		if !ok {
			cs = make(map[string]bool)
			componentStates[component] = cs
		}
		cs[state] = true
	}

	// check which components we need to add
	for _, comp := range m.r.Data.StreamConfig.Streams {
		regCmpState(comp.ComponentId, comp.StateId)

		cmp, err := ct.CreateComponentIfNotExists(comp.ComponentId)
		if err != nil {
			return err
		}
		st, err := cmp.CreateStateIfNotExists(comp.StateId, comp.Config)
		if err != nil {
			return err
		}

		// local component + state
		lcmp, err := lct.CreateComponentIfNotExists(comp.ComponentId)
		if err != nil {
			return err
		}
		lst, err := lcmp.CreateStateIfNotExists(comp.StateId, comp.Config)
		if err != nil {
			return err
		}

		// Backfill
		if err := st.Backfill(lst); err != nil {
			return err
		}
		lst.RemoteStates[m.r.Data.Id] = st
	}

	// check which components we need to remove
	for _, cmp := range ct.getAllComponents() {
		stateIds, ok := componentStates[cmp.Name]
		// if we have the component, maybe we don't have one of the states
		if ok {
			for _, state := range cmp.States {
				if _, ok := stateIds[state.Name]; ok {
					continue
				}
				m.unregisterRemote(cmp.Name, state.Name)
				if err := cmp.DeleteState(state.Name); err != nil {
					return err
				}
			}
			continue
		}
		// Delete the entire component
		for _, state := range cmp.States {
			m.unregisterRemote(cmp.Name, state.Name)
		}
		if err := ct.DeleteComponent(cmp.Name); err != nil {
			return err
		}
	}

	return nil
}

func (m *RemoteManager) unregisterRemote(component, state string) {
	// local component + state
	lct := m.r.RemoteList.reporter.LocalTree
	lcmp, err := lct.GetComponent(component)
	if err != nil || lcmp == nil {
		return
	}
	lst, err := lcmp.GetState(state)
	if err != nil || lst == nil {
		return
	}
	delete(lst.RemoteStates, m.r.Data.Id)
}

// Push all states to remote, return nil if successful.
func (m *RemoteManager) pushStateToRemote(st *State) error {
	// Right off the bat, if nothing is in the stream, exit
	if len(st.Data.AllTimestamp) == 0 {
		return nil
	}

	// Sanity check
	if st.Data.LastTimestamp == 0 {
		st.Data.LastTimestamp = st.Data.AllTimestamp[len(st.Data.AllTimestamp)-1]
	}

	sendIdx := 0

	// find next index to send
	for st.Data.AllTimestamp[sendIdx] <= st.Data.RemoteTimestamp {
		sendIdx++
		if sendIdx == len(st.Data.AllTimestamp) {
			return nil
		}
	}

	// Until RemoteTimestamp >= LastTimestamp
	for st.Data.RemoteTimestamp < st.Data.LastTimestamp {
		nextSend := st.Data.AllTimestamp[sendIdx]
		nextEntry, err := st.getEntry(nextSend)
		if err != nil {
			return err
		}

		// Serialize to json
		jsonData, err := json.Marshal(nextEntry.Data)
		if err != nil {
			return err
		}

		// Send the entry
		resp, err := m.client.PushStreamEntry(context.Background(), &remote.PushStreamEntryRequest{
			Context: &remote.RequestContext{
				HostIdentifier: m.r.RemoteList.reporter.hostIdentifier,
				ComponentId:    st.Component.Name,
				StateId:        st.Name,
			},
			Entry: &remote.RemoteStreamEntry{
				Timestamp: nextSend,
				EntryType: int32(nextEntry.Type),
				JsonData:  string(jsonData),
			},
			ConfigCrc32: m.r.Data.StreamConfig.Crc32,
		})

		if err != nil {
			// we had some kind of push error
			glog.Warningf("Error pushing entry to remote, reconnecting: %v\n", err)

			// disconnect
			m.conn.Close()
			m.client = nil
			m.conn = nil

			return err
		}

		if resp.Config != nil {
			m.r.Data.StreamConfig = resp.Config
			m.r.WriteToDb() // ignore error here, should be handled later
		}

		// Update the new timing
		st.Data.RemoteTimestamp = nextSend
		st.WriteToDb() //  any errors here should be picked up later

		// Always keep at least 1 snapshot + everything past it
		// Once we pass a snapshot, delete everything before it
		if nextEntry.Type == stream.StreamEntrySnapshot && sendIdx != 0 {
			if err := st.PurgeBefore(sendIdx); err != nil {
				// ignore error here, shouldn't affect this process.
				glog.Warningf("Error purging remote entries, %v", err)
			}
			// index 0 will now be sendIdx
			sendIdx = 0
		}

		sendIdx++

	}
	return nil
}

// Push all stream entries to remote
func (m *RemoteManager) pushToRemote() {
	for state := range m.pendingPush {
		if err := m.pushStateToRemote(state); err != nil {
			return
		}
		delete(m.pendingPush, state)
	}
	m.flushDirtyChan(true)
}

func (m *RemoteManager) update() bool {
	m.flushDirtyChan(true)
	if m.client == nil {
		glog.Infof("Attempting to connect to %s...", m.r.Data.Config.Endpoint)
		if err := m.connect(); err != nil {
			glog.Warningf("Failed to connect to remote, %v.", err)
			return !m.shouldStopTimeout(time.Duration(5) * time.Second)
		}
		glog.Infof("Connected to %s.", m.r.Data.Config.Endpoint)
		// force a config re-fetch
		m.r.Data.StreamConfig = nil
	}
	// Do this frequently to avoid a block
	m.flushDirtyChan(true)
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
		if err := m.syncRemoteConfig(); err != nil {
			glog.Warningf("Error when syncing remote config, %v - continuing, but beware.", err)
		}
	}

	if m.pendingPush == nil {
		m.buildPushSlice()
	}

	// Push as much as we can to the server
	m.pushToRemote()

	// This will happen only when we have an error.
	if len(m.pendingPush) != 0 {
		return true
	}

	// Re-check the server if nothing happens for a minute or so.
	minTick := time.NewTimer(time.Duration(1) * time.Minute)
	select {
	case <-m.stopChan:
		return false
	case state := <-m.r.StateDirtyChan:
		if _, ok := m.pendingPush[state]; ok {
			m.pendingPush[state]++
		} else {
			m.pendingPush[state] = 1
		}
		m.flushDirtyChan(true)
	case <-minTick.C:
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
