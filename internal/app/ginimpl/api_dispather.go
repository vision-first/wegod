package ginimpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/enum"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/auth"
	"net/http"
	"time"
)

type ApiContext struct {
	auth *auth.Auth
	ginCtx *gin.Context
}

func NewApiContext(ctx *gin.Context) *ApiContext {
	authForInterface, ok := ctx.Get(enum.GinCtxKeyAuth)
	apiCtx := &ApiContext{
		ginCtx: ctx,
	}
	if ok {
		apiCtx.auth = authForInterface.(*auth.Auth)
	}
	return apiCtx
}

func (c *ApiContext) GetAuth() (*auth.Auth, bool) {
	if c.auth == nil {
		return nil, false
	}

	return c.auth, true
}

func (c *ApiContext) GetGinCtx() *gin.Context {
	return c.ginCtx
}

func (c *ApiContext) Err() error {
	return c.ginCtx.Err()
}

func (c *ApiContext) Deadline() (deadline time.Time, ok bool) {
	return c.Deadline()
}

func (c *ApiContext) Done() <-chan struct{} {
	return c.ginCtx.Done()
}

func (c *ApiContext) Value(key interface{}) interface{} {
	return c.Value(key)
}

var _ api.Context = (*ApiContext)(nil)

type ApiDispatcher struct {
	dispatcher *api.Dispatcher
}

func NewApiDispatcher() *ApiDispatcher {
	return &ApiDispatcher{
		dispatcher: api.NewDispatcher(facades.MustLogger()),
	}
}

func (d *ApiDispatcher) Dispatch(ctx *gin.Context, apiMethod interface{}) {
	var dispatchApiReqHandler api.DispatchApiReqFunc
	switch ctx.Request.Method {
	case http.MethodPost:
		dispatchApiReqHandler = parsePostCtxToMakeDispatchApiReqHandler(ctx)
	case http.MethodGet:
		dispatchApiReqHandler = parseGetCtxToMakeDispatchApiReqHandler(ctx)
	default:
		err := errors.New("not support http method")
		facades.MustLogger().Warn(ctx, err)
		err = ctx.AbortWithError(http.StatusForbidden, err)
		if err != nil {
			facades.MustLogger().Error(ctx, err)
		}
		return
	}
	resp, err := d.dispatcher.Dispatch(NewApiContext(ctx), apiMethod, dispatchApiReqHandler)
	if err != nil {
		facades.MustLogger().Error(ctx, err)
		response.EndFailedJsonWithErr(ctx, err)
		return
	}
	response.EndSuccessfulJson(ctx, resp)
}

func parseGetCtxToMakeDispatchApiReqHandler(ctx *gin.Context) api.DispatchApiReqFunc {
	return func(req interface{}) error {
		if err := ctx.ShouldBindQuery(req); err != nil {
			facades.MustLogger().Error(ctx, err)
		}
		return nil
	}
}

func parsePostCtxToMakeDispatchApiReqHandler(ctx *gin.Context) api.DispatchApiReqFunc {
	return func(req interface{}) error {
		var err error
		switch ctx.ContentType() {
		case "application/json":
			err = ctx.ShouldBindJSON(req)
		default:
			err = errors.New("not support http request content type")
		}
		if err != nil {
			return err
		}
		return nil
	}
}