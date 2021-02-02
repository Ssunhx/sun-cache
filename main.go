package main

import (
	"fmt"
	"log"
	"net/http"
	"sun-cache/cache"
	hp "sun-cache/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	cache.NewGroup("score", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowdb] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%v not exist", key)
		}))

	addr := "localhost:9099"
	peers := hp.NewHTTPPool(addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
