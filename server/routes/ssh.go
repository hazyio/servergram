package routes

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
)

func registerSSHRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/ssh", api, database)
	server.Post("/add", http.StatusCreated, addSSHKey)

}

func addSSHKey(data internal.NormalRouteData[AddSSHKeyRequest]) (*EmptyResponse, error) {
	input := data.Input
	name := input.Body.Name
	database := data.Database
	key := input.Body.Key
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
	encryptedSSHKey, encryptedSSHKeyErr := internal.EncryptString(key, input.Body.Key)
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

	return &EmptyResponse{}, nil
}
