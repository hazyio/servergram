package sshmanager

import "sync"

type SSHManager struct {
	mu sync.Mutex
	
}

func NewSSHManager() *SSHManager {
	return &SSHManager{}
}

func (m *SSHManager) ConnectNew() {
	// Implementation for connecting a new SSH session goes here
}
