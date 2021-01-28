package cache

import (
	lru2 "sun-cache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex // 互斥锁
	lru        *lru2.Cache
	cacheBytes int64
}

// 先判断 lru 是否为 nil， 如果等于nil就进行实例化，称为延迟初始化 lazy initialization，一个对象的延迟初始化意味着对象
// 的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求。
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru2.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		// 断言，将 v 转换成 ByteView 类型
		return v.(ByteView), ok
	}
	return
}
