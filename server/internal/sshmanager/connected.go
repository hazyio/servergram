package sshmanager

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hazyio.com/servergram/server/agents"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal/common"
	"hazyio.com/servergram/shared"
)

type SSHConnectedData struct {
	Id     uint
	Name   string
	client *ssh.Client
	Agent  *Agent
	OS     db.OS
	Arch   db.Arch
	Home   string
}

func newUnEncryptedConnection(server *db.Server) (*SSHConnectedData, error) {
	ctx := context.Background()
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
	os, arch, err := getOSAndArch(ctx, client)
	if err != nil {
		return nil, err
	}
	agentBin, err := agents.GetAgentFor(os, arch)
	if err != nil {
		return nil, err
	}
	home, err := getRemoteHome(ctx, client)
	if err != nil {
		return nil, err
	}
	uploadCtx, uploadCancel := context.WithTimeout(ctx, 200*time.Second)
	defer uploadCancel()

	remotePath, err := common.RunWithContext(uploadCtx, func() (string, error) {
		return uploadAgent(home, client, agentBin)
	}, nil)
	if err != nil {
		slog.Error("Failed to upload agent", "error", err)
		return nil, err
	}

	agentConn, err := startAgentConnection(ctx, client, remotePath)
	if err != nil {
		return nil, err
	}

	return &SSHConnectedData{Id: server.ID, Name: server.Name, client: client, OS: os, Arch: arch, Home: home, Agent: NewAgent(agentConn)}, nil
}
func (c *SSHConnectedData) IsAlive() bool {
	client := c.client // read once; avoids racing with a reconnect swapping the field
	if client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// TODO: Consider implementing a more robust keepalive mechanism to detect dead connections.
	_, err := common.RunWithContext(ctx, func() (any, error) {
		_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
		return nil, err
	}, func(any) { client.Close() })
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

func uploadAgent(home string, c *ssh.Client, bin []byte) (string, error) {
	agentPath := home + "/.servergram/agent"
	checksum := fmt.Sprintf("%x", sha256.Sum256(bin))
	checkSession, err := c.NewSession()
	if err != nil {
		return "", err
	}
	remoteChecksum, checkErr := checkSession.Output(`if [ -f "$HOME/.servergram/agent" ]; then sha256sum "$HOME/.servergram/agent"; fi`)
	_ = checkSession.Close()
	if checkErr == nil {
		fields := strings.Fields(string(remoteChecksum))
		if len(fields) > 0 && fields[0] == checksum {
			slog.Info("Agent is up to date", "home", home)
			return agentPath, nil
		}
	}

	slog.Info("Uploading agent", "home", home)
	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	dir := fmt.Sprintf("%s/.servergram/", home)
	cmd := fmt.Sprintf(
		`set -e; mkdir -p %[1]s; cat > %[1]s/agent.new; chmod 700 %[1]s/agent.new; mv -f %[1]s/agent.new %[1]s/agent`, dir)

	// create a std input pipe for the session
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}
	var stderr bytes.Buffer
	session.Stderr = &stderr

	if err := session.Start(cmd); err != nil {
		// returns immediately, command is running
		return "", err
	}

	if _, err := io.Copy(stdin, bytes.NewReader(bin)); err != nil {
		return "", fmt.Errorf("upload: %w (remote: %s)", err, stderr.String())
	}
	if err := stdin.Close(); err != nil {
		// EOF: lets cat finish
		return "", err
	}
	if err := session.Wait(); err != nil {
		// *ssh.ExitError on non-zero exit
		return "", fmt.Errorf("install: %w (remote: %s)", err, stderr.String())
	}
	return agentPath, nil
}
func startAgentConnection(ctx context.Context, c *ssh.Client, remotePath string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	slog.Info("Starting agent connection", "remotePath", remotePath)
	session, sessionErr := c.NewSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	return common.RunWithContext(ctx, func() (*grpc.ClientConn, error) {

		stdin, _ := session.StdinPipe()
		stdout, _ := session.StdoutPipe()
		// sess.Stderr =slog.Default().Handler()

		if err := session.Start(remotePath); err != nil {
			session.Close()
			return nil, err
		}

		conn := shared.NewStdioConn(stdout, stdin, func() error {
			stdin.Close()
			return session.Close()
		})

		slog.Info("Agent connection established")
		var used atomic.Bool
		return grpc.NewClient("passthrough:///agent",
			grpc.WithTransportCredentials(insecure.NewCredentials()), // SSH already encrypts and authenticates
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				if used.Swap(true) {
					session.Close()
					return nil, errors.New("agent connection is gone")
				}
				return conn, nil
			}),
		)
	}, func(conn *grpc.ClientConn) {
		conn.Close()
		session.Close() // kills the remote process's stdin/channel
	})

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
func archFromUname(s string) (db.Arch, error) {
	switch strings.TrimSpace(s) {
	case "x86_64", "amd64":
		return db.ArchAMD64, nil
	case "aarch64", "arm64": // Linux reports aarch64, macOS reports arm64
		return db.ArchARM64, nil
	default:
		return "", fmt.Errorf("unsupported architecture %q", s)
	}
}
func getOSAndArch(ctx context.Context, c *ssh.Client) (db.OS, db.Arch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	sess, err := c.NewSession()
	if err != nil {
		return "", "", err
	}
	defer sess.Close()
	type osInfo struct {
		OS   db.OS
		Arch db.Arch
	}

	info, err := common.RunWithContext(ctx, func() (osInfo, error) {
		var b bytes.Buffer
		sess.Stdout = &b
		if err := sess.Run("uname -s && uname -m"); err != nil {
			return osInfo{}, err
		}
		lines := strings.Split(strings.TrimSpace(b.String()), "\n")
		if len(lines) != 2 {
			return osInfo{}, errors.New("unexpected uname output")
		}
		arch, err := archFromUname(lines[1])
		if err != nil {
			return osInfo{}, err
		}
		return osInfo{db.OS(strings.ToLower(lines[0])), arch}, nil
	}, func(_ osInfo) { sess.Close() })
	if err != nil {
		return "", "", err
	}
	return info.OS, info.Arch, nil

}

func getRemoteHome(ctx context.Context, c *ssh.Client) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	return common.RunWithContext(ctx, func() (string, error) {

		defer session.Close()
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run("echo $HOME"); err != nil {
			return "", err
		}
		return strings.TrimSpace(b.String()), nil
	}, func(home string) {
		// Cleanup function for RunWithContext: close the SSH session if it is still open
		session.Close()
	})

}
