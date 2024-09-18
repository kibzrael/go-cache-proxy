package main

import (
	"fmt"
	"kibzrael/cacheproxy/cmd/cacheproxy"
	"os"
	"strconv"
)

func main(){
	args := os.Args[1:]

	var port uint64
	var origin string
	clearCache := false

	for i, arg := range args{
		if arg == "--port"{
			p, err := strconv.ParseUint(args[i+1], 10, 64)
			if err != nil{
				panic(err)
			}
			port = p
		} else if arg == "--origin"{
			origin = args[i+1]
		} else if arg == "--clear-cache"{
			clearCache = true
		}
	}

	if clearCache {
		cacheproxy.ClearCache()
	} else if port == 0 {
		fmt.Println("Port argument required --port")
	} else if origin == "" {
		fmt.Println("Origin argument required --origin")
	} else {
		cacheproxy.Proxy(port, origin)
	}
}