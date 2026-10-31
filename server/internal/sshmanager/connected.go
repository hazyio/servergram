package sshmanager

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"hazyio.com/servergram/db"
)

type SSHConnectedData struct {
	Id     uint
	Name   string
	client *ssh.Client
}

func newUnEncryptedConnection(server *db.Server) (*SSHConnectedData, error) {
	signer, err := ssh.ParsePrivateKey([]byte(server.SSHKey.Key))
	if err != nil {
		slog.Error("failed to parse private key", "error", err)
		return nil, err
	}
	var hostErr error
	cfg := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: dbHostKeyCallback(server, &hostErr),
		Timeout:         10 * time.Second,
	}
	slog.Info("attempting to connect to SSH server", "host", server.Host, "port", server.Port)
	addr := net.JoinHostPort(server.Host, strconv.Itoa(int(server.Port)))

	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		// ssh.Dial flattens the callback error to a string, so return the original typed error.
		if hostErr != nil {
			return nil, hostErr
		}
		slog.Error("failed to dial SSH", "error", err)
		return nil, err
	}

	return &SSHConnectedData{Id: server.ID, Name: server.Name, client: client}, nil
}
func (c *SSHConnectedData) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

type UnknownHostError struct{ Host, Fingerprint, KeyType string }

func (e *UnknownHostError) Error() string {
	return fmt.Sprintf("unknown host key for %s:%s", e.Host, e.Fingerprint)
}

type HostKeyMismatchError struct{ Host, Expected, Got string }

func (e *HostKeyMismatchError) Error() string {
	return fmt.Sprintf("host key mismatch for %s: expected %s, got %s", e.Host, e.Expected, e.Got)
}

func dbHostKeyCallback(server *db.Server, captured *error) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		got := ssh.FingerprintSHA256(key)

		if !server.Trusted || server.HostKeyFingerprint == "" {
			*captured = &UnknownHostError{Host: hostname, Fingerprint: got, KeyType: key.Type()}
			return *captured
		}
		if server.HostKeyFingerprint != got {
			*captured = &HostKeyMismatchError{Host: hostname, Expected: server.HostKeyFingerprint, Got: got}
			return *captured
		}
		return nil
	}
}
