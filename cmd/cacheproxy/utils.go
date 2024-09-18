package cacheproxy

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/patrickmn/go-cache"
)

func ForwardResponse(res http.ResponseWriter, response *http.Response) {
	headers := res.Header()

	for key, value := range headers {
		for _, v := range value {
			res.Header().Set(key, v)
		}
	}
	res.WriteHeader(response.StatusCode)
	io.Copy(res, response.Body)
}

func CheckCache(request string, c *cache.Cache, res http.ResponseWriter, req *http.Request) bool {
	body, found := c.Get(request)
	if found {
		log.Println("Cache found")
		r := bufio.NewReader(bytes.NewReader(body.([]byte)))
		cacheResponse, err := http.ReadResponse(r, req)
		if err != nil {
			log.Println("ERROR: ", err)
		} else {
			res.Header().Set("X-Cache", "HIT")
			ForwardResponse(res, cacheResponse)
			return true
		}
	}
	return false
}

func CacheResponse(request string, c *cache.Cache, response *http.Response) {
	b, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	c.Set(request, b, cache.NoExpiration)
}
