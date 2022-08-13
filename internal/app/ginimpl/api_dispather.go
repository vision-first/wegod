package ginimpl

import (
	"errors"
	"github.com/995933447/apperrdef"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/enums"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/auth"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"net/http"
	"time"
)

type ApiContext struct {
	ginCtx *gin.Context
	ident *auth.Ident
}

var _ api.Context = (*ApiContext)(nil)

func NewApiContext(ctx *gin.Context) *ApiContext {
	authForInterface, ok := ctx.Get(enums.GinCtxKeyAuth)
	apiCtx := &ApiContext{
		ginCtx: ctx,
	}
	if ok {
		apiCtx.ident = authForInterface.(*auth.Ident)
	}
	return apiCtx
}

func (c *ApiContext) GetAuthIdent() (*auth.Ident, bool, error) {
	if c.ident == nil {
		return nil, false, nil
	}
	return c.ident, true, nil
}

func (c *ApiContext) GetAuthIdentOrFailed() (*auth.Ident, error) {
	if c.ident == nil {
		return nil, apperrdef.NewError(errs.ErrCodeUnauthorized)
	}
	return c.ident, nil
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

func (d *ApiDispatcher) MakeDispatchFunc(apiMethod interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		d.Dispatch(ctx, apiMethod)
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