package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/go-redis/redis/v8"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/encrypt"
	"github.com/vision-first/wegod/internal/pkg/enums"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"time"
)

type Auth struct {
	logger *log.Logger
}

func NewAuth(logger *log.Logger) *Auth {
	return &Auth{
		logger: logger,
	}
}

func (a *Auth) RememberPhoneVerifyCodeForAuth(ctx context.Context, phone, code string) error {
	err := facades.RedisGroup(a.logger).Set(ctx, enums.MakeRedisKeyPhoneVerifyCode(phone), []byte(code), time.Minute * 20)
	if err != nil {
		a.logger.Error(ctx, err)
		return a.TransErr(err)
	}
	return nil
}

func (a *Auth) AuthPhoneVerifyCode(ctx context.Context, phone, code string) (bool, error) {
	codeBytes, err := facades.RedisGroup(a.logger).Get(ctx, enums.MakeVerifyCodeMsgForRegister(phone))
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		a.logger.Error(ctx, err)
		return false, a.TransErr(err)
	}

	return string(codeBytes) == code, nil
}

func (a *Auth) CreateAuthTokenForUser(ctx context.Context, userId uint64) (string, error) {
	token, err := encrypt.NewJWT(config.Conf.Encrypt.SigningKey).GenerateToken(userId)
	if err != nil {
		a.logger.Error(ctx, err)
		return "", a.TransErr(err)
	}
	return token, nil
}

func (a *Auth) TransErr(err error) error {
	return err
}