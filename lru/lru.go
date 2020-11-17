package lru

import (
	"container/list"
)

// LRU cache
type Cache struct {
	maxbytes int64                    // 允许使用的最大内存
	nbytes   int64                    // 当前已经使用的内存
	ll       *list.List               // go 语言中标准库实现的双向链表 list.List
	cache    map[string]*list.Element // map，键是字符串，值是双向链表中对应节点的指针

	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil
}

// 双向链表节点的数据类型，在链表中仍然保存每个值对应的key的好处在于，淘汰队首节点时，需要用key从字典删除映射
type entry struct {
	key   string
	value Value
}

// 允许实现了 Value 接口的任意类型，
type Value interface {
	Len() int
}

// 初始化cache
func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxbytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// 查找 从字典中找到对应的双向链表的节点，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 链表节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，这里约定front为队尾）
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除操作，实际上就是缓存淘汰，移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	// 取到队首节点，从链表中删除
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从字典中删除节点的映射关系
		delete(c.cache, kv.key)
		// 更新当前内存占用大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 若有回调函数，执行回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// 键存在，更新对应节点的值，并将节点移动到队尾。
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 键不存在是新增的场景，首先队尾添加新节点，并在字典添加key和节点的映射关系
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 超过了最大值，则删除最近最少访问的节点
	for c.maxbytes != 0 && c.maxbytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 获取添加了多少数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
