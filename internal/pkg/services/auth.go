package services

import (
	"context"

	"github.com/995933447/log-go"
	"github.com/995933447/stringhelper-go"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

type Auth struct {
	logger *log.Logger
}

func (a *Auth) IssueATempCertificate(ctx context.Context) string {
	tempCertificate := stringhelper.GenRandomStr(stringhelper.RandomStringModNumberPlusLetter, 32)
	facades.RedisGroup(a.logger).Set(ctx, )
}
