package cache

import "log"


//按照类型创建实例
func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache(ttl)
}
	if typ == "rocksdb" {
		c = newRocksdbCache(ttl)
	}
	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
