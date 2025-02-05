package initializers

import (
	"github.com/NurymGM/jwt-token/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
