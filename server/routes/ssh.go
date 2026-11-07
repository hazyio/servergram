package routes

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/routes/models"
)

func registerSSHRoutes(api *huma.Group, database *gorm.DB, allowedOrigins []string, polls *internal.Polls) {
	server := internal.NewGroup("/ssh-key", api, database, []string{"SSH"}, nil, allowedOrigins, polls)
	server.Post("/add", http.StatusCreated, addSSHKey)
	server.Get("/list", http.StatusOK, listSSHKeys)
	server.Get("/view-key/{id}", http.StatusOK, viewSSHKey)
	server.Delete("/delete/{id}", http.StatusNoContent, deleteSSHKey)
}

func addSSHKey(data internal.RouteData[models.AddSSHKeyRequest]) (*EmptyRR, error) {
	input := data.Input
	name := input.Body.Name
	database := data.Database
	key := input.Body.Key
	env := internal.LoadEnv()
	err := input.Resolve()
	if err != nil {
		return nil, err
	}
	find := database.Where("name = ?", name).First(&db.SSHKey{})
	if find.Error == nil {
		return nil, internal.NewHumaError(http.StatusConflict, "SSH key with that name already exists")
	}
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, internal.InternalServerError
	}
	encryptedSSHKey, encryptedSSHKeyErr := internal.EncryptString(key, env.EncryptionMasterKey)
	if encryptedSSHKeyErr != nil {
		slog.Error("An error occurred while encrypting the SSH key", "error", encryptedSSHKeyErr)
		return nil, internal.InternalServerError
	}
	newKey := db.SSHKey{
		Name:      name,
		Key:       encryptedSSHKey,
		Encrypted: input.Encrypted,
	}
	if err := database.Create(&newKey).Error; err != nil {
		slog.Error("An error occurred while saving ssh key", "error", err)
		return nil, internal.InternalServerError
	}
	db.AddLog(db.LogActionCreate, "SSH key \""+name+"\" added", data.Database)
	return &EmptyRR{}, nil
}

func listSSHKeys(data internal.RouteData[EmptyRR]) (*models.SSHKeyListResponse, error) {
	database := data.Database
	var total int64
	if err := database.Model(&db.SSHKey{}).Count(&total).Error; err != nil {
		slog.Error("count SSH keys failed", "error", err)
		return nil, internal.InternalServerError
	}
	response := &models.SSHKeyListResponse{}
	response.Body.Items = make([]models.SSHKeyListItem, 0)
	database.Model(&db.SSHKey{}).
		Select("ssh_keys.id, ssh_keys.name, ssh_keys.created_at, ssh_keys.encrypted, COUNT(servers.id) AS resources_attached").
		Joins("LEFT JOIN servers ON servers.ssh_key_id = ssh_keys.id").
		Group("ssh_keys.id").
		Scan(&response.Body.Items)

	return response, nil
}

func viewSSHKey(data internal.RouteData[PathIdRequest]) (*models.ViewSSHKeyResponse, error) {
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
	decryptedKey, err := internal.DecryptString(key.Key, env.EncryptionMasterKey)
	if err != nil {
		slog.Error("An error occurred while decrypting SSH key", "error", err)
		return nil, internal.InternalServerError
	}
	response := &models.ViewSSHKeyResponse{}
	response.Body.Private = decryptedKey
	db.AddLog(db.LogActionUpdate, "SSH key \""+key.Name+"\" viewed", data.Database)
	return response, nil
}

func deleteSSHKey(data internal.RouteData[PathIdRequest]) (*EmptyRR, error) {
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
	db.AddLog(db.LogActionDelete, "SSH key \""+key.Name+"\" deleted", data.Database)
	return &EmptyRR{}, nil
}
