package handler

import (
	"count/backend/internal/model"
	"count/backend/pkg/auth"

	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func New(db *gorm.DB, jwtSecret string) *Handler {
	return &Handler{
		DB:        db,
		JWTSecret: []byte(jwtSecret),
	}
}

type authPayload struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

type listResponse struct {
	Items any `json:"items"`
}

func tokenForUser(secret []byte, user model.User) (string, error) {
	return auth.GenerateToken(secret, user.ID, user.Role, user.StationRole, user.Phone, user.Username)
}
