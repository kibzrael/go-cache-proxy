package cacheproxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func CacheRequest(origin string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Fetching %s from %s", req.URL.Path, origin)
		response, err := http.Get(origin + req.URL.Path)
		if err != nil{
			panic(err)
		}
		
		if response.StatusCode != 200{
			fmt.Fprintln(res, "Request Failed with status code: ", response.StatusCode)
			return;
		}

		headers := res.Header()

		for key , value := range headers{
			for _, v := range value{
				res.Header().Set(key, v)
			}
		}
		res.WriteHeader(response.StatusCode)
		io.Copy(res, response.Body)
	}
}