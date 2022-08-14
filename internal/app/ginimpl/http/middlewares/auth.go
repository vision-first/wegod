package middlewares

import (
	"github.com/995933447/apperrdef"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/encrypt"
	"github.com/vision-first/wegod/internal/pkg/errs"
)

func Auth(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Token")

	if len(token) == 0 {
		// 未登录或非法登录
		response.EndFailedJsonWithErr(ctx, apperrdef.NewErr(errs.ErrCodeUnauthorized))
		return
	}

	_, err, expired := encrypt.NewJWT(config.Conf.Encrypt.Jwt.SigningKey).ParseToken(token)
	if err != nil {
		facades.MustLogger().Error(ctx, err)
		response.EndFailedJsonWithErr(ctx, apperrdef.NewErr(errs.ErrCodeUnauthorized))
		return
	}

	if expired {
		facades.MustLogger().Error(ctx, err)
		response.EndFailedJsonWithErr(ctx, apperrdef.NewErrWithMsg(errs.ErrCodeUnauthorized, "登录已过期"))
		return
	}

	ctx.Next()
}
