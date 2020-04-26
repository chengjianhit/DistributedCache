package lrustruct

import (
	"container/list"
	"sync"
)

type CacheFacade struct {
	mu         sync.Mutex
	lru        *Cache
	cacheBytes int64
}

type Cache struct {
	maxSize   int64      //最大内存
	usedSize  int64      //已使用内存大小
	ll        *list.List //双向链表
	cacheMap  map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func (c *Cache) Len() (size int) {
	return c.ll.Len()
}

func NewCache(maxSize int64, onEnvicted func(string, Value)) *Cache {
	return &Cache{
		maxSize:   maxSize,
		ll:        list.New(),
		cacheMap:  make(map[string]*list.Element),
		OnEvicted: onEnvicted,
	}
}

func (cacheFacade *CacheFacade) add(key string, value BaseStore) {
	cacheFacade.mu.Lock()
	defer cacheFacade.mu.Unlock()
	//初始化延迟
	if cacheFacade.lru == nil {
		cacheFacade.lru = NewCache(cacheFacade.cacheBytes, nil)
	}
	cacheFacade.lru.Set(key, value)
}

func (cacheFacade *CacheFacade) getValue(key string) (value BaseStore, ok bool) {
	cacheFacade.mu.Lock()
	if cacheFacade.lru == nil {
		return
	}

	if v, succ := cacheFacade.lru.Get(key); succ {
		return v.(BaseStore), succ
	}

	return
}

func (c *Cache) Set(key string, value Value) {
	if existedVal, ok := c.cacheMap[key]; ok {
		c.ll.MoveToFront(existedVal) //最新访问的ele，放到队列头部
		entry := existedVal.Value.(*entry)
		c.usedSize += int64(value.Len() - entry.value.Len())
		entry.value = value
	} else {
		newEntry := c.ll.PushFront(&entry{key, value})
		c.cacheMap[key] = newEntry
		c.usedSize += int64(len(key)) + int64(value.Len())
	}

	for c.maxSize != -1 && c.maxSize < c.usedSize {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, succ bool) {
	if element, ok := c.cacheMap[key]; ok {
		c.ll.MoveToFront(element)
		entry := element.Value.(*entry)
		return entry.value, true
	}

	return nil, false

}

func (c *Cache) RemoveOldest() {
	lastEle := c.ll.Back()
	if lastEle != nil {
		c.ll.Remove(lastEle)
		entry := lastEle.Value.(*entry)
		delete(c.cacheMap, entry.key)
		c.usedSize -= int64(len(entry.key)) + int64(entry.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}
}
