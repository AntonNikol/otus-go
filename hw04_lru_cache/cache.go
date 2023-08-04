package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cacheItem struct {
	key   Key
	value interface{}
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var cachedValue interface{}

	node, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(node)
		item := c.getNodeItem(node)
		cachedValue = item.value
	}

	return cachedValue, ok
}

func (c *lruCache) getNodeItem(node *ListItem) *cacheItem {
	return node.Value.(*cacheItem)
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.items[key]

	if ok {
		item := c.getNodeItem(node)
		item.value = value
		c.queue.MoveToFront(node)
		return ok
	}

	if c.queue.Len() == c.capacity {
		c.removeLastFromQueue()
	}
	item := newCacheItem(key, value)
	node = c.queue.PushFront(item)
	c.items[key] = node
	return ok
}

func (c *lruCache) removeLastFromQueue() {
	leastUsedNode := c.queue.Back()
	leastUsedItem := c.getNodeItem(leastUsedNode)
	delete(c.items, leastUsedItem.key)
	c.queue.Remove(leastUsedNode)
}
