package cacheproxy

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

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

	file, _ := CacheFile()
	defer file.Close()
	writer := CacheWriter{
		file: file,
	}
	c.Save(writer)
}

const FILE_NAME string = "Cache"

func CacheFile() (*os.File, bool) {
	file, err := os.OpenFile(FILE_NAME, os.O_RDWR, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create(FILE_NAME)
			if err != nil {
				panic(err)
			}
			return file, true
		} else {
			panic(err)
		}
	}
	return file, false
}

func LoadCache(c *cache.Cache) {
	file, _ := CacheFile()
	defer file.Close()
	reader := CacheReader{
		file: file,
	}
	c.Load(reader)
}

type CacheWriter struct {
	file *os.File
}

func (w CacheWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

type CacheReader struct {
	file *os.File
}

func (r CacheReader) Read(b []byte) (int, error) {
	return r.file.Read(b)
}
