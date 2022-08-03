package api

import (
	"context"
	"github.com/vision-first/wegod/internal/pkg/auth"
)

type Context interface {
	context.Context
	GetAuth() (*auth.Auth, bool)
}