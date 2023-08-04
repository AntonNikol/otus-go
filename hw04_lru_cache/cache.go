package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		value := item.Value
		fmt.Println(value)
		return item.Value, true
	}
	return nil, false
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}

	if c.queue.Len() >= c.capacity {
		old := c.queue.Back()
		if old != nil {
			delete(c.items, old.Value.(Key))
			c.queue.Remove(old)
		}
	}

	item := c.queue.PushFront(value)
	c.items[key] = item

	return false
}
