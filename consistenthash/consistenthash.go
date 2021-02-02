package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash 函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // key 虚拟节点的哈希值， value 真是节点名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		// 默认的 hash 算法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 传入 0 个或多个真是节点，每个节点对应一个真实节点 key， 对应的创建 m.replicas 个虚拟节点，虚拟节点的名称的
// strconv.Itoa(i) + key, 通过添加编号的方式区分不同的虚拟节点，使用 m.hash() 计算虚拟节点的哈希值，并且添加到
// m.keys 哈希环上，在 m.hashMap 中增加虚拟节点和真实节点的映射关系。最后对哈希环上的值排序
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 一个真实节点映射 replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 计算 key 的 hash 值
	hash := int(m.hash([]byte(key)))

	// 顺时针查找第一个节点的 index， 因为 m.keys 是个环，所以取余数
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
