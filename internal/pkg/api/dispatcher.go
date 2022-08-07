package api

import (
	"errors"
	"github.com/995933447/log-go"
	"reflect"
)

type DispatchApiReqFunc func(interface{}) error

type Dispatcher struct {
	logger *log.Logger
}

func NewDispatcher(logger *log.Logger) *Dispatcher {
	return &Dispatcher{
		logger: logger,
	}
}

func (d *Dispatcher) Dispatch(ctx Context, apiMethod interface{}, dispatchApiReq DispatchApiReqFunc) (interface{}, error) {
	apiMethodType := reflect.TypeOf(apiMethod)
	invalidApiMethodErr := errors.New("arg[1] is not func(api.*Context, interface{}) (interface{}, error)")
	if apiMethodType.Kind() != reflect.Func {
		err := errors.New("arg[1] is not a func")
		d.logger.Error(ctx, err)
		return nil, err
	}
	if apiMethodType.NumIn() != 2 || apiMethodType.NumOut() != 2 {
		d.logger.Error(ctx, invalidApiMethodErr)
		return nil, invalidApiMethodErr
	}
	if _, ok := reflect.New(apiMethodType.In(0)).Interface().(*Context); !ok {
		d.logger.Error(ctx, invalidApiMethodErr)
		return nil, invalidApiMethodErr
	}
	apiReqType := apiMethodType.In(1)
	apiReqKing := apiReqType.Kind()
	if apiReqKing == reflect.Ptr {
		apiReqType = apiReqType.Elem()
	}
	if apiReqType.Kind() != reflect.Struct {
		err := errors.New("arg[1].In(1) is not a struct")
		d.logger.Error(ctx, err)
		return nil, err
	}
	apiRespType := apiMethodType.Out(0)
	if apiRespType.Kind() == reflect.Ptr {
		apiRespType = apiRespType.Elem()
	}
	if apiRespType.Kind() != reflect.Struct {
		err := errors.New("arg[1].Out(1) is not a struct")
		d.logger.Error(ctx, err)
		return nil, err
	}
	if _, ok := reflect.New(apiMethodType.Out(1)).Interface().(*error); ok {
		d.logger.Error(ctx, invalidApiMethodErr)
		return nil, invalidApiMethodErr
	}
	apiReq := reflect.New(apiReqType)
	if err := dispatchApiReq(apiReq.Interface()); err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}
	apiReplyTypes := reflect.ValueOf(apiMethod).Call([]reflect.Value{reflect.ValueOf(ctx), apiReq})
	err := apiReplyTypes[1].Interface().(error)
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}
	resp := apiReplyTypes[0].Interface()
	if resp == nil {
		resp = reflect.New(apiRespType)
	}
	return resp, nil
}
