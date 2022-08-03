package auth

import (
	"github.com/vision-first/wegod/internal/pkg/datamodels"
)

type Ident struct {
	user *datamodels.User
}

func NewIdent(user *datamodels.User) *Ident {
	return &Ident{
		user: user,
	}
}

func (i *Ident) GetUser() *datamodels.User {
	return i.user
}
