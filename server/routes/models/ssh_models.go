package models

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/ssh"
)

type AddSSHKeyRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" maxLength:"100" pattern:"^[A-Za-z0-9_ .@-]+$" doc:"The name of the SSH key"`
		Key  string `json:"key" required:"true" maxLength:"16384" doc:"The private SSH key (may be passphrase-encrypted)"`
	}
	Encrypted bool
}
type ViewSSHKeyResponse struct {
	Body struct {
		Private string `json:"private"`
	}
}
type SSHKeyListItem struct {
	Id                uint      `json:"id"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
	Encrypted         bool      `json:"encrypted"`
	ResourcesAttached uint      `json:"resources_attached"`
}
type SSHKeyListResponse struct {
	Body struct {
		Items []SSHKeyListItem `json:"items" required:"true" doc:"The list of SSH keys"`
	}
}

func (r *AddSSHKeyRequest) Resolve() error {
	// normalize (mutates the request, so the handler sees the cleaned value)
	key := strings.ReplaceAll(strings.TrimSpace(r.Body.Key), "\r\n", "\n") + "\n"
	r.Body.Key = key
	r.Body.Name = strings.TrimSpace(r.Body.Name)

	fail := func(msg string) error {
		return &huma.ErrorDetail{
			Location: "body.key",
			Message:  msg,
		}
	}

	if r.Body.Name == "" {
		return &huma.ErrorDetail{Location: "body.name", Message: "name is required"}
	}
	// reject public keys with a helpful message
	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key)); err == nil {
		return fail("that looks like a public key, paste the private key")
	}

	_, err := ssh.ParsePrivateKey([]byte(key))
	if err == nil {
		return nil // valid, unencrypted
	}
	var missing *ssh.PassphraseMissingError
	if errors.As(err, &missing) {
		slog.Info("private key is passphrase-protected")
		r.Encrypted = true
		return nil // valid, but passphrase-protected
	}

	return fail("not a valid private key")
}
