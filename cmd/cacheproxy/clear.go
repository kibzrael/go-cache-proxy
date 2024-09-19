package cacheproxy

import "os"

func ClearCache() {
	os.Remove(FILE_NAME)
}
