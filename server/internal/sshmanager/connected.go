package sshmanager

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/shared"
)

type SSHConnectedData struct {
	Id     uint
	Name   string
	client *ssh.Client
	// Agent  *grpc.ClientConn
	OS   db.OS
	Arch db.Arch
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
	os, arch, err := getOSAndArch(client)
	if err != nil {
		return nil, err
	}
	return &SSHConnectedData{Id: server.ID, Name: server.Name, client: client, OS: os, Arch: arch}, nil
}
func (c *SSHConnectedData) IsAlive() bool {
	if c.client == nil {
		return false
	}
	// TODO: Consider implementing a more robust keepalive mechanism to detect dead connections.
	_, _, err := c.client.SendRequest("keepalive@openssh.com", true, nil)
	return err == nil
}
func (c *SSHConnectedData) Disconnect() error {
	if c.client != nil {
		if !c.IsAlive() {
			return nil
		}
		return c.client.Close()
	}
	return nil
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

func getOSAndArch(c *ssh.Client) (db.OS, db.Arch, error) {
	sess, err := c.NewSession()
	if err != nil {
		return "", "", err
	}
	defer sess.Close()

	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run("uname -s && uname -m"); err != nil {
		return "", "", err
	}
	output := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(output) != 2 {
		return "", "", errors.New("unexpected uname output")
	}
	return db.OS(strings.ToLower(output[0])), db.Arch(output[1]), nil
}

// func uploadAgent(c *ssh.Client) (string, error) {
// 	sess, err := c.NewSession()
// 	if err != nil {
// 		return "", err
// 	}
// 	return "sjsj", nil
// }

func dialAgent(c *ssh.Client, remotePath string) (*grpc.ClientConn, error) {
	sess, err := c.NewSession()
	if err != nil {
		return nil, err
	}

	stdin, _ := sess.StdinPipe()
	stdout, _ := sess.StdoutPipe()
	// sess.Stderr =slog.Default().Handler()

	if err := sess.Start(remotePath); err != nil {
		return nil, err
	}

	conn := shared.NewStdioConn(stdout, stdin, func() error {
		stdin.Close()
		return sess.Close()
	})

	var used atomic.Bool
	return grpc.NewClient("passthrough:///agent",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // SSH already encrypts and authenticates
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			if used.Swap(true) {
				return nil, errors.New("agent connection is gone")
			}
			return conn, nil
		}),
	)
}
