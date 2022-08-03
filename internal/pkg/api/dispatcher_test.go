package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/apidispatcher"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/pkg/api/apis"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"reflect"
	"testing"
)

func TestReflectApiMethod(t *testing.T) {
	buddhaApi := apis.NewBuddha(facades.MustLogger())
	reflectType := reflect.TypeOf(buddhaApi.PageBuddha)
	if reflectType.Kind() != reflect.Func {
		t.Fatalf("not func, reflect type is:%s", reflectType.Kind().String())
	}
	for i := 0; i < reflectType.NumIn(); i++ {
		if reflectType.In(i).Kind() == reflect.Interface {
			if _, ok := reflect.New(reflectType.In(i)).Interface().(*Context); ok {
				t.Logf("field %d:%v", i, reflectType.In(i).String())
			} else {
				t.Logf("ty:%t", reflect.New(reflectType.In(i)).Interface())
			}
			continue
		}

		t.Logf("field %d:%v", i, reflectType.In(i).Elem())
	}

	apiReplyReflects := reflect.ValueOf(buddhaApi.PageBuddha).Call([]reflect.Value{reflect.ValueOf(apidispatcher.NewApiContext(&gin.Context{})), reflect.ValueOf(&dtos.PageBuddhaReq{})})
	err := apiReplyReflects[1].Interface().(error)
	if err != nil {
		t.Logf("err:%v", err)
	}
	resp := apiReplyReflects[0].Interface()
	t.Logf("err:%v", resp)

	arg1val := reflect.New(reflectType.In(1).Elem())
	t.Logf("arg1val:%+v", arg1val.Interface())
	j := `{"limit":10}`
	err = json.Unmarshal([]byte(j), arg1val.Interface())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", arg1val)
}
