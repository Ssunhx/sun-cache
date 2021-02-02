package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sun-cache/cache"
)

const defaultBasePath = "/_suncache/"

type HTTpPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTpPool {
	return &HTTpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTpPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check url path
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serveing unexpected path" + r.URL.Path)
	}

	p.Log("%s, %s", r.Method, r.URL.Path)

	// <basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupname, key := parts[0], parts[1]
	// 根据 groupname 获取对应的 group
	group := cache.GetGroup(groupname)

	if group == nil {
		http.Error(w, "error group name", http.StatusNotFound)
		return
	}
	// 根据 key 在 group 中查找
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, "key not exist", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
}
