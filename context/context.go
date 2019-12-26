package context

import (
	goContext "context"
	"encoding/hex"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type CancelFunc goContext.CancelFunc

// context的本意是希望可以参考golang官方context包
// 并能额外提供一些更加便捷的操作方式
type Context struct {
	sync.RWMutex

	ctx    goContext.Context
	parent *Context
	name   string
	where  string
	logger *zap.SugaredLogger

	track  *uint64
	values map[interface{}]interface{}

	L *zap.SugaredLogger
}

var (
	_ goContext.Context = &Context{}

	bootid  string = uuid.NewV4().String()
	counter uint64 = 0
)

// sync boot uuid
func BootID(id string) {
	u, err := uuid.FromString(id)
	if err != nil {
		panic(err)
	}
	bootid = u.String()
}

// New创建一个新的随机名称Context，当前位置at可选
func New(at ...string) *Context {
	seq := atomic.AddUint64(&counter, 1)
	var b = []byte{
		byte(seq >> 40),
		byte(seq >> 32),
		byte(seq >> 24),
		byte(seq >> 16),
		byte(seq >> 8),
		byte(seq),
	}
	var buf = make([]byte, 2*len(b))
	hex.Encode(buf, b)
	return NewWithName(bootid[:24]+string(buf), at...)
}

// NewWithName创建一个指定名称的Context，当前位置at可选
func NewWithName(name string, at ...string) *Context {
	var t uint64 = 0
	var w string
	if len(at) > 0 && at[0] != "" {
		w = at[0]
	}
	logger := zap.S().Named(name)
	c := &Context{
		ctx:    goContext.Background(),
		parent: nil,
		name:   name,
		where:  w,
		logger: logger,

		track:  &t,
		values: map[interface{}]interface{}{},

		L: logger,
	}
	if c.where != "" {
		c.L = c.L.With("@", c.where)
	}
	return c
}

// ~=copy
func (c *Context) shadow() *Context {
	return &Context{
		L: c.L,

		ctx:    c.ctx,
		parent: c,
		name:   c.name,
		where:  c.where,
		logger: c.logger,

		track:  c.track,
		values: map[interface{}]interface{}{},
	}
}

// Context派生子Context，用于并发派生新goroutine的场景
// 比如内部api基本上都采用了并发调用的方式，需要一个额外的追踪标记来标明日志的从属请求
func (c *Context) New(at ...string) *Context {
	shadow := c.shadow()

	if len(at) > 0 && at[0] != "" {
		if shadow.where == "" {
			shadow.where = at[0]
		} else {
			shadow.where += "." + at[0]
		}
	}

	// 保留ctx，应用新name
	seq := atomic.AddUint64(shadow.track, 1)
	var t uint64 = 0

	shadow.logger = shadow.logger.Named(fmt.Sprintf("%d", seq))
	shadow.L = shadow.logger
	if shadow.where != "" {
		shadow.L = shadow.L.With("@", shadow.where)
	}
	shadow.track = &t

	return shadow
}

// 用于标记当前位置，或者传达调用路径时使用
func (c *Context) At(at string) *Context {
	shadow := c.shadow()
	if at != "" {
		if shadow.where == "" {
			shadow.where = at
		} else {
			shadow.where += "." + at
		}
	}
	if shadow.where != "" {
		shadow.L = shadow.logger.With("@", shadow.where)
	}
	return shadow
}

func (c *Context) Name() string {
	return c.name
}

func (c *Context) Session() string {
	return c.name
}

func (c *Context) WithCancel() (*Context, CancelFunc) {
	shadow := c.shadow()

	// 保留name，应用新ctx
	ctx, f := goContext.WithCancel(shadow.ctx)
	shadow.ctx = ctx

	return shadow, CancelFunc(f)
}

func (c *Context) WithDeadline(d time.Time) (*Context, CancelFunc) {
	shadow := c.shadow()

	// 保留name，应用新ctx
	ctx, f := goContext.WithDeadline(shadow.ctx, d)
	shadow.ctx = ctx

	return shadow, CancelFunc(f)
}

func (c *Context) WithTimeout(timeout time.Duration) (*Context, CancelFunc) {
	shadow := c.shadow()

	// 保留session，应用新ctx
	ctx, f := goContext.WithTimeout(shadow.ctx, timeout)
	shadow.ctx = ctx

	return shadow, CancelFunc(f)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	val, _ := c.Get(key)
	return val
}

func (c *Context) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	val, ok := c.values[key]
	c.RUnlock()
	if !ok && c.parent != nil {
		return c.parent.Get(key)
	}
	return val, ok
}

func (c *Context) GetString(key interface{}) string {
	val, ok := c.Get(key)
	if ok {
		return val.(string)
	}
	return ""
}

func (c *Context) GetInt(key interface{}) int {
	val, ok := c.Get(key)
	if ok {
		return val.(int)
	}
	return 0
}

func (c *Context) GetInt64(key interface{}) int64 {
	val, ok := c.Get(key)
	if ok {
		return val.(int64)
	}
	return 0
}

func (c *Context) GetUint(key interface{}) uint {
	val, ok := c.Get(key)
	if ok {
		return val.(uint)
	}
	return 0
}

func (c *Context) GetUint64(key interface{}) uint64 {
	val, ok := c.Get(key)
	if ok {
		return val.(uint64)
	}
	return 0
}

func (c *Context) Set(key interface{}, value interface{}) {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	c.Lock()
	c.values[key] = value
	c.Unlock()
}

func (c *Context) Del(key interface{}) {
	c.Lock()
	delete(c.values, key)
	c.Unlock()
}
