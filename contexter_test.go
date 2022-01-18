package contexter_test

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"testing"

	"github.com/KarpelesLab/contexter"
)

//go:noinline
func getTest(ctx context.Context, s string) context.Context {
	// NOTE it is important that the value is being used at least once
	ctx.Value(nil)
	return getTest2()
}

func getTest2() context.Context {
	return contexter.Context()
}

//go:noinline
func getTestAny(ctx context.Context, s string) context.Context {
	// NOTE it is important that the value is being used at least once
	ctx.Value(nil)
	return getTestAny2()
}

func getTestAny2() context.Context {
	var ctx context.Context
	if !contexter.Find(&ctx) {
		log.Printf("NOT FOUND")
	}
	return ctx
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	log.Printf("ctx = %p", ctx)

	ctx2 := getTest(ctx, "hello world")
	log.Printf("ctx2 = %p", ctx2)

	if ctx != ctx2 {
		t.Errorf("invalid value returned: %p", ctx2)
	}

	ctx3 := getTestAny(ctx, "hello world")
	log.Printf("ctx3 = %p", ctx3)

	if ctx != ctx3 {
		t.Errorf("invalid value returned in any: %p", ctx3)
	}
}

type TestObj struct{}

func (t *TestObj) MarshalJSON() ([]byte, error) {
	ctx := contexter.Context()
	if ctx == nil {
		return nil, errors.New("could not fetch context")
	}

	res := map[string]interface{}{"foo": ctx.Value("test")}
	return json.Marshal(res)
}

//go:noinline
func encodeJson(ctx context.Context, obj interface{}) ([]byte, error) {
	res, err := json.Marshal(obj)
	runtime.KeepAlive(ctx)
	return res, err
}

func TestJson(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "bar")
	obj := &TestObj{}

	val, err := encodeJson(ctx, obj)

	if err != nil {
		t.Errorf("json test failed: %s", err)
		return
	}

	if string(val) != `{"foo":"bar"}` {
		t.Errorf("json output failed, should be {\"foo\":\"bar\"} but got %s", val)
	}
}
