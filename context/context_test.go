package context

import (
	gcontext "context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func _() { fmt.Print(); zap.S() }

func init() {
	//zap.ReplaceGlobals(zap.NewExample())
}

func TestContext(t *testing.T) {
	{
		ctx := New(gcontext.Background())
		ctx.L.Infow("hello", "name", "CJey")
		ctx1 := ctx.New()
		ctx1.L.Infow("hello", "name", "CJey")
		ctx2 := ctx.New()
		ctx2.L.Infow("hello", "name", "CJey")
		ctx21 := ctx2.New()
		ctx21.L.Infow("hello", "name", "CJey")
		ctx22 := ctx2.New()
		ctx22.L.Infow("hello", "name", "CJey")

		ctx22a := ctx22.WithTimeout(time.Second)
		ctx22a.L.Infow("hello", "name", "CJey")

		ctx22b := ctx22.At("lname")
		ctx22b.L.Infow("hello", "name", "CJey")

		ctx22ba := ctx22b.At("nname")
		ctx22ba.L.Infow("hello", "name", "CJey")
	}
	{
		ctx := New(gcontext.Background())
		ctx.L.Infow("hello", "name", "CJey")
		ctx1 := ctx.New()
		ctx1.L.Infow("hello", "name", "CJey")
		ctx2 := ctx.New()
		ctx2.L.Infow("hello", "name", "CJey")
		ctx21 := ctx2.New()
		ctx21.L.Infow("hello", "name", "CJey")
		ctx22 := ctx2.New()
		ctx22.L.Infow("hello", "name", "CJey")

		ctx22a := ctx22.WithTimeout(time.Second)
		ctx22a.L.Infow("hello", "name", "CJey")

		ctx22b := ctx22.At("lname")
		ctx22b.L.Infow("hello", "name", "CJey")

		ctx22ba := ctx22b.At("nname")
		ctx22ba.L.Infow("hello", "name", "CJey")
	}
}
