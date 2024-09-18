package cacheproxy

import (
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
)

func Proxy(port uint64, origin string, c *cache.Cache){
	server := http.Server{
		Addr: fmt.Sprintf(":%v", port),
	}

	http.HandleFunc("/", CacheRequest(origin, c))

	fmt.Printf("Application listening on http://localhost:%v\n", port)
	err := server.ListenAndServe()
	if err != nil{
		panic(err)
	}
}