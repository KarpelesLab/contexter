package contexter_test

import (
	"context"
	"log"
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
