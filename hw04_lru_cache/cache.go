package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity uint64
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItem struct {
	value interface{}
	key   Key
}

func NewCache(capacity uint64) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		item.Value = cacheItem{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() == c.capacity {
		if toRemove := c.queue.Back(); toRemove != nil {
			if removeItem, ok := toRemove.Value.(cacheItem); ok {
				delete(c.items, removeItem.key)
			}
			c.queue.Remove(toRemove)
		}
	}

	c.items[key] = c.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		if value, ok := item.Value.(cacheItem); ok {
			return value.value, true
		}
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
