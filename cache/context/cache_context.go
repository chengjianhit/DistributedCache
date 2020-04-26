package context

import (
	"CacheAdmin/cache/lrustruct"
	"sync"
)

type Loader interface {
	load(key string) ([]byte, error)
}

type LoaderFunc func(key string) ([]byte, error)

func (loaderFunc LoaderFunc) Load(key string) ([]byte, error) {
	return loaderFunc(key)
}

type CacheContext struct {
	name        string
	cacheFacade lrustruct.CacheFacade
	loadFunc    Loader
}

var (
	muLock       sync.RWMutex
	contextGroup = make(map[string]*CacheContext)
)

func NewCacheContext(name string, cacheBytes int64, loader Loader) {

}
