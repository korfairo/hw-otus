package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// cacheItem is an element of the cache inside *Node.Value.
type cacheItem struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Node

	mu sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Node, capacity),
		mu:       sync.RWMutex{},
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, exists := l.items[key]
	if exists {
		l.updateValue(key, value, node)
		return true
	}
	l.addValue(key, value)
	l.checkOverflow()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	node, exists := l.items[key]
	if !exists {
		return nil, false
	}

	l.queue.MoveToFront(node)

	if node.Value != nil {
		return node.Value.(cacheItem).Value, true
	}

	return nil, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*Node, l.capacity)
}

func (l *lruCache) updateValue(key Key, newValue interface{}, node *Node) {
	l.queue.Remove(node)
	l.addValue(key, newValue)
}

func (l *lruCache) addValue(key Key, value interface{}) {
	newNode := l.queue.PushFront(cacheItem{
		Key:   key,
		Value: value,
	})
	l.items[key] = newNode
}

func (l *lruCache) checkOverflow() {
	if l.queue.Len() > l.capacity {
		tail := l.queue.Back()
		l.queue.Remove(tail)
		delete(l.items, tail.Value.(cacheItem).Key)
	}
}
