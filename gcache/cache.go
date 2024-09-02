package gcache

import (
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Load 读取数据
	// key 缓存key必须全局唯一
	// duration 缓存有效时间
	// ds 当缓存没有数据时从数据源读取
	Load(key any, ds func() (out any, err error), duration time.Duration) (out any, err error)

	Set(key, value any, duration time.Duration)
	Get(key any) (value any)
}

// Load 加载数据，如果缓存里有则从缓存读取
func Load[T any](c Cache, key any, ds func() (any, error), duration time.Duration) (out T, err error) {
	r, err := c.Load(key, ds, duration)
	if err != nil {
		return
	}
	out = r.(T)
	return
}

func Set(c Cache, k, v any, duration time.Duration) {
	c.Set(k, v, duration)
}
func Get[T any](c Cache, k any) (v T) {
	val := c.Get(k)
	if val == nil {
		return
	}
	v = val.(T)
	return
}
