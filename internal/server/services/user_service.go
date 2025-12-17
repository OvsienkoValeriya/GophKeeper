package services

import (
	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/OvsienkoValeriya/GophKeeper/internal/server/auth"
)

func NewUser(username, password string) (*models.User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Username: username,
		Password: hash,
	}, nil
}
