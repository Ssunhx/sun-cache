package cache

import (
	"fmt"
	"log"
	"sync"
)

// Getter 和 回调函数
type Getter interface {
	Get(key string) ([]byte, error)
}

// 函数类型 GetterFunc， 并且实现了 Getter 接口的 Get 方法
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 一个 Group 认为是一个命名空间，每个 Group 拥有一个唯一的名称 name。
// getter Getter 缓存命中时获取数据的回调 callback
// mainCache cache 一开始实现的并发缓存
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 实例化 Group， 并且将 group 存储在全局变量 groups中
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

// 获取特定名称的 group，使用了只读锁
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 从 mainCache 中查找缓存，如果存在就返回
// 如果不存在，就调用 load 方法， load 调用 getLocally （分布式场景下嗲用 getFromPeer 从其他节点获取），
// getLocally 调用用户回调函数 g.getter.Get() 获取源数据，并且将源数据添加到缓存 mainCache 中
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Fatalln("hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 使用回掉函数来获取 key 对应的 value
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 将 key-value 添加到缓存中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
