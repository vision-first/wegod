package auth

import (
	"github.com/vision-first/wegod/internal/pkg/datamodels"
)

type Auth struct {
	user *datamodels.User
}

func NewAuth(user *datamodels.User) *Auth {
	return &Auth{
		user: user,
	}
}

func (a *Auth) User() *datamodels.User {
	return a.user
}
