package timeLRU

import (
	"github.com/golang/groupcache/lru"
	"sync"
	"time"
)

type TimeLRU struct {
	cache *lru.Cache
	sync.RWMutex
}

type Element struct {
	value interface{}
	expireTime time.Time
}

func New(max int) *TimeLRU {
	return &TimeLRU{cache:lru.New(max)}
}

func (lru *TimeLRU)Add(key string, value interface{}, duration time.Duration)  {
	lru.Lock()
	defer lru.Unlock()
	lru.cache.Add(key, Element{
		value:      value,
		expireTime: time.Now().Add(duration),
	})
}

func (lru *TimeLRU)Get(key string) (value interface{}, ok bool) {
	lru.RLock()
	var elementS interface{}
	elementS, ok = lru.cache.Get(key)
	lru.RUnlock()
	if !ok {
		return
	}
	element := elementS.(Element)
	if element.expireTime.Sub(time.Now()) < 0 {
		lru.Remove(key)
		return nil, false
	}
	return element.value, ok
}

func (lru *TimeLRU)Remove(key string) {
	lru.Lock()
	defer lru.Unlock()
	lru.cache.Remove(key)
}

func (lru *TimeLRU)RemoveOldest() {
	lru.Lock()
	defer lru.Unlock()
	lru.cache.RemoveOldest()
}

func (lru *TimeLRU)Clear() {
	lru.Lock()
	defer lru.Unlock()
	lru.cache.Clear()
}