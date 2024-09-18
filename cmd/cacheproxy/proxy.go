package cacheproxy

import (
	"fmt"
	"net/http"
)

func Proxy(port uint64, origin string){
	server := http.Server{
		Addr: fmt.Sprintf(":%v", port),
	}

	http.HandleFunc("/", CacheRequest(origin))

	fmt.Printf("Application listening on http://localhost:%v\n", port)
	err := server.ListenAndServe()
	if err != nil{
		panic(err)
	}
}