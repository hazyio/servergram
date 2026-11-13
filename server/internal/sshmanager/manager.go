package sshmanager

import (
	"errors"
	"log/slog"
	"sync"

	"hazyio.com/servergram/server/db"
)

var (
	SSHManagerErrHostTrustNeeded = errors.New("host trust needed")
)

type SSHConnectedListData struct {
	Name string
	Id   uint
}
type SSHManager struct {
	mu        sync.Mutex
	connected map[uint]*SSHConnectedData
	queued    map[uint]SSHConnectedListData
}

func NewSSHManager() *SSHManager {
	return &SSHManager{
		connected: make(map[uint]*SSHConnectedData),
		queued:    make(map[uint]SSHConnectedListData),
	}
}

func (m *SSHManager) IsQueued(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.queued[id]
	return ok
}
func (m *SSHManager) Queue(id uint, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queued[id] = SSHConnectedListData{Id: id, Name: name}
}
func (m *SSHManager) GetQueue() []SSHConnectedListData {
	m.mu.Lock()
	defer m.mu.Unlock()
	data := make([]SSHConnectedListData, 0)
	for _, v := range m.queued {
		data = append(data, v)
	}
	return data
}
func (m *SSHManager) Dequeue(id uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.queued, id)
}
func (m *SSHManager) Disconnect(id uint) error {
	m.mu.Lock()

	connection, ok := m.connected[id]
	if !ok {
		m.mu.Unlock()
		return nil
	}
	if err := connection.Disconnect(); err != nil {
		m.mu.Unlock()
		return err
	}
	
	delete(m.connected, id)
	m.mu.Unlock()
	slog.Info("SSH connection closed", "server_id", id)
	return nil
}

func (m *SSHManager) IsConnected(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.connected[id]
	if ok && !data.IsAlive() {
		ok = false
	}
	return ok
}
func (m *SSHManager) GetConnected(id uint) (*SSHConnectedData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	connection, ok := m.connected[id]
	if !ok {
		return nil, errors.New("connection not found")
	}
	if !connection.IsAlive() {
		delete(m.connected, id)
		return nil, errors.New("connection not alive")
	}
	return connection, nil
}
func (m *SSHManager) GetAllConnected() []SSHConnectedListData {
	m.mu.Lock()
	defer m.mu.Unlock()
	data := make([]SSHConnectedListData, 0)
	for id, sshKey := range m.connected {
		data = append(data, SSHConnectedListData{
			Id:   id,
			Name: sshKey.Name,
		})
	}
	return data
}

func (m *SSHManager) ConnectNewUnEncrypted(server *db.Server) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	connection, err := newUnEncryptedConnection(server, func(err error) {
		m.Disconnect(server.ID)
	})
	if err != nil {
		return err
	}
	m.connected[server.ID] = connection
	// remove from queuedIds after connecting
	delete(m.queued, server.ID)
	slog.Info("SSH connection established", "server_id", server.ID, "server_name", server.Name)
	return nil
}
