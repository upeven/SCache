package cache

import (
	"sync"
	"time"
)
//声明值 结构体，包含值和创建时间
type value struct {
	v       []byte
	created time.Time
}

type inMemoryCache struct {
	//底层数据结构，一个map，key为string,value为value结构体
	c     map[string]value
	//锁
	mutex sync.RWMutex
	Stat
	ttl time.Duration
}

//新增键值
func (c *inMemoryCache) Set(k string, v []byte) error {
	//加锁。并发安全
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = value{v, time.Now()}
	c.add(k, v)
	return nil
}

//获取键值
func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k].v, nil
}

//删除键值
func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v.v)
	}
	return nil
}

//获取缓存服务状态
func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

//创建一个自动扫描键值对过期时间的实例
func newInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{make(map[string]value), sync.RWMutex{}, Stat{}, time.Duration(ttl) * time.Second}
	if ttl > 0 {
		go c.expirer()
	}
	return c
}

//删除过期缓存并释放空间
func (c *inMemoryCache) expirer() {
	for {
		time.Sleep(c.ttl)
		//添加读锁
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}
