package routes

import (
	"errors"
	"net/http"
	"time"

	"strings"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/ssh"
)

type EmptyRR struct {
}
type LoginRequest struct {
	Body struct {
		Password string `json:"password" required:"true" doc:"The password for login" minLength:"16" maxLength:"128"`
	}
	UserAgent string `header:"User-Agent" required:"true" hidden:"true"`
}

type LoginSuccessfulRedirectResponse struct {
	Status    int
	SetCookie http.Cookie `header:"Set-Cookie"`
}
type MessageResponse struct {
	Message string `json:"message"`
}

type AddSSHKeyRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" maxLength:"100" pattern:"^[A-Za-z0-9_ .@-]+$" doc:"The name of the SSH key"`
		Key  string `json:"key" required:"true" maxLength:"16384" doc:"The private SSH key (may be passphrase-encrypted)"`
	}
}
type CreateGroupRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" maxLength:"100" pattern:"^[A-Za-z0-9_ .@-]+$" doc:"The name of the group"`
	}
}

type SearchRequest struct {
	Body struct {
		Query string `json:"query" required:"true" doc:"The search query"`
	}
}
type SearchResponse[T any] struct {
	Body struct {
		Items []T `json:"items" required:"true" doc:"The search results"`
		Total int `json:"total" required:"true" doc:"The total number of search results"`
	}
}
type SSHKeySearchData struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Encrypted bool   `json:"encrypted"`
}

type GroupSearchData struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type GroupListResponse struct {
	Body struct {
		Items []struct {
			Id                uint      `json:"id"`
			Name              string    `json:"name"`
			CreatedAt         time.Time `json:"created_at"`
			ResourcesAttached uint      `json:"resources_attached"`
		} `json:"items" required:"true" doc:"The list of groups"`
		Total int `json:"total" required:"true" doc:"The total number of groups"`
	}
}
type SSHKeyListResponse struct {
	Body struct {
		Items []struct {
			Id                uint      `json:"id"`
			Name              string    `json:"name"`
			CreatedAt         time.Time `json:"created_at"`
			Encrypted         bool      `json:"encrypted"`
			ResourcesAttached uint      `json:"resources_attached"`
		} `json:"items" required:"true" doc:"The list of SSH keys"`
		Total int `json:"total" required:"true" doc:"The total number of SSH keys"`
	}
}
type PathIdRequest struct {
	ID uint `path:"id" minimum:"1" doc:"The ID of the current resource"`
}
type ViewSSHKeyResponse struct {
	Body struct {
		Private string `json:"private"`
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

	return fail("not a valid private key")
}

// Ensure you call [Resolve()] before calling [KeyEncrypted()] to normalize and validate the key.
func (r *AddSSHKeyRequest) KeyEncrypted() bool {
	_, err := ssh.ParsePrivateKey([]byte(r.Body.Key))
	if err == nil {
		return false
	}

	var missing *ssh.PassphraseMissingError
	if errors.As(err, &missing) {
		return true // valid structure, needs a passphrase
	}

	return false // not a valid private key at all meaning it's invalid

}
