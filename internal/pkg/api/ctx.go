package api

import (
	"context"
	"github.com/vision-first/wegod/internal/pkg/auth"
)

type Context interface {
	context.Context
	GetAuthIdent() (*auth.Ident, bool, error)
	GetAuthIdentOrFailed() (*auth.Ident, error)
}