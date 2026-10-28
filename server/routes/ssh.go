package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
)

func registerSSHRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/ssh", api, database)
	server.Post("/add", http.StatusCreated, addSSHKey)
	server.Get("/list", http.StatusOK, listSSHKeys)
	server.Get("/view-key/{id}", http.StatusOK, viewSSHKey)
	server.Delete("/delete/{id}", http.StatusNoContent, deleteSSHKey)

}

func addSSHKey(data internal.NormalRouteData[AddSSHKeyRequest]) (*EmptyRR, error) {
	input := data.Input
	name := input.Body.Name
	database := data.Database
	key := input.Body.Key
	env := internal.LoadEnv()
	err := input.Resolve()
	if err != nil {
		return nil, err
	}
	encrypted := input.KeyEncrypted()
	find := database.Where("name = ?", name).First(&db.SSHKey{})
	if find.Error == nil {
		return nil, internal.NewHumaError(http.StatusConflict, "SSH key with that name already exists")
	}
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, internal.InternalServerError
	}
	encryptedSSHKey, encryptedSSHKeyErr := internal.EncryptString(key, env.ShhMasterKey)
	if encryptedSSHKeyErr != nil {
		slog.Error("An error occurred while encrypting the SSH key", "error", encryptedSSHKeyErr)
		return nil, internal.InternalServerError
	}
	newKey := db.SSHKey{
		Name:      name,
		Key:       encryptedSSHKey,
		Encrypted: encrypted,
	}
	if err := database.Create(&newKey).Error; err != nil {
		slog.Error("An error occurred while saving ssh key", "error", err)
		return nil, internal.InternalServerError
	}

	return &EmptyRR{}, nil
}

func listSSHKeys(data internal.NormalRouteData[EmptyRR]) (*SSHKeyListResponse, error) {
	database := data.Database
	var keys []db.SSHKey
	if err := database.Find(&keys).Error; err != nil {
		slog.Error("An error occurred while fetching SSH keys", "error", err)
		return nil, internal.InternalServerError
	}
	response := &SSHKeyListResponse{}
	response.Body.Items = make([]struct {
		Id                uint      `json:"id"`
		Name              string    `json:"name"`
		CreatedAt         time.Time `json:"created_at"`
		Encrypted         bool      `json:"encrypted"`
		ResourcesAttached uint      `json:"resources_attached"`
	}, len(keys))
	database.Model(&db.SSHKey{}).
		Select("ssh_keys.id, ssh_keys.name, ssh_keys.created_at, COUNT(servers.id) AS resources_attached").
		Joins("LEFT JOIN servers ON servers.ssh_key_id = ssh_keys.id").
		Group("ssh_keys.id").
		Scan(&response.Body.Items)

	response.Body.Total = len(keys)

	return response, nil
}

func viewSSHKey(data internal.NormalRouteData[PathIdRequest]) (*ViewSSHKeyResponse, error) {
	env := internal.LoadEnv()

	database := data.Database
	id := data.Input.ID
	var key db.SSHKey
	if err := database.First(&key, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, internal.NewHumaError(http.StatusNotFound, "SSH key not found")
		}
		slog.Error("An error occurred while fetching SSH key", "error", err)
		return nil, internal.InternalServerError
	}
	decryptedKey, err := internal.DecryptString(key.Key, env.ShhMasterKey)
	if err != nil {
		slog.Error("An error occurred while decrypting SSH key", "error", err)
		return nil, internal.InternalServerError
	}
	response := &ViewSSHKeyResponse{}
	response.Body.Private = decryptedKey

	return response, nil
}

func deleteSSHKey(data internal.NormalRouteData[PathIdRequest]) (*EmptyRR, error) {
	database := data.Database
	id := data.Input.ID
	var key db.SSHKey
	if err := database.First(&key, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, internal.NewHumaError(http.StatusNotFound, "SSH key not found")
		}
		slog.Error("An error occurred while fetching SSH key", "error", err)
		return nil, internal.InternalServerError
	}

	if err := database.Delete(&db.SSHKey{}, id).Error; err != nil {
		slog.Error("An error occurred while deleting SSH key", "error", err)
		return nil, internal.InternalServerError
	}
	// Check if there are any servers attached to the SSH key
	var count int64
	if err := database.Model(&db.Server{}).Where("ssh_key_id = ?", id).Count(&count).Error; err != nil {
		slog.Error("An error occurred while checking attached servers", "error", err)
		return nil, internal.InternalServerError
	}
	if count > 0 {
		return nil, internal.NewHumaError(http.StatusConflict, "Cannot delete SSH key with attached servers")
	}
	if err := database.Delete(&key).Error; err != nil {
		slog.Error("An error occurred while deleting SSH key", "error", err)
		return nil, internal.InternalServerError
	}
	return &EmptyRR{}, nil
}
