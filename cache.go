package smart

import (
	"github.com/astaxie/beego/cache"
	cache2 "github.com/goburrow/cache"
	"sync"
	"tianwei.pro/beego-guava"
	"time"
)

type CacheManager interface {
	Get(name string) cache.Cache
}

type SmartCacheManager struct {
	sync.Mutex
	container map[string]cache.Cache
}

func NewSmartCacheManager() CacheManager{
	return &SmartCacheManager{
		container: make(map[string]cache.Cache),
	}
}

func (s *SmartCacheManager) Get(name string) cache.Cache {
	if v, ok := s.container[name]; ok {
		return v
	} else {
		s.Lock()
		defer s.Unlock()
		if v, ok := s.container[name]; ok {
			return v
		}


		c := cache2.NewLoadingCache(func(key cache2.Key) (value cache2.Value, e error) {
			return nil, nil
		}, cache2.WithMaximumSize(1000),
			cache2.WithExpireAfterAccess(30 * time.Hour),)
		cc := beego_guava.NewGuava(c)
		s.container[name] = cc
		return cc
	}
}