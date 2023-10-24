package gcache

import "time"

// Cache 缓存接口
type Cache interface {
	// Load 读取数据
	// key 缓存key必须全局唯一
	// duration 缓存有效时间
	// ds 当缓存没有数据时从数据源读取
	Load(key string, duration time.Duration, ds func() (out any, err error)) (out any, err error)
}

// Load 加载数据，如果缓存里有则从缓存读取
func Load[T any](c Cache, key string, duration time.Duration, ds func() (any, error)) (out T, err error) {
	r, err := c.Load(key, duration, ds)
	if err != nil {
		return
	}
	out = r.(T)
	return
}
