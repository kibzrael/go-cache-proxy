package cacheproxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"
)

func CacheRequest(origin string, c *cache.Cache) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Fetching %s from %s", req.URL.Path, origin)
		request := origin + req.URL.Path

		found := CheckCache(request, c, res, req)
		if found {
			return
		}
		res.Header().Set("X-Cache", "MISS")

		response, err := http.Get(request)
		if err != nil {
			log.Println("ERROR: ", err)
			fmt.Fprintln(res, "Request Failed: ", err)
			return
		}

		if response.StatusCode != 200 {
			fmt.Fprintln(res, "Request Failed with status code: ", response.StatusCode)
			return
		}

		bodyBytes, _ := io.ReadAll(response.Body)
		// Restore the io.ReadCloser to its original state
		response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		ForwardResponse(res, response)
		response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		CacheResponse(request, c, response)
	}
}
