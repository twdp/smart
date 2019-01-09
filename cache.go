package smart

import "github.com/astaxie/beego/cache"

type CacheManager interface {
	
}

type SmartCacheManager struct {
	container map[string]cache.Cache
}

func NewSmartCacheManager() CacheManager{
	return &SmartCacheManager{
		container: make(map[string]cache.Cache),
	}
}