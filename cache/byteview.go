package cache

// b 存储真实的缓存值，选择 byte 类型是为了支持任意类型的数据，字符串、图片等
type ByteView struct {
	b []byte
}

// 在 lru.Cache 的实现中，要求被缓存对象必须实现 Value 接口，即 Len（）方法，返回占用内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// b 是只读的，返回一个拷贝，防止缓存值在外部程序被修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
