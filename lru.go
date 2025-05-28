package main

import (
	"net/http"
	"sync"
)

type CacheEntry struct {
	StatusCode int
	Headers    http.Header
	Data       []byte
}

type Node struct {
	url   string
	value CacheEntry
	next  *Node
	prev  *Node
}

func NewNode(url string, value CacheEntry) *Node {
	return &Node{
		url:   url,
		value: value,
	}
}

type LRUCache struct {
	capacity int
	latest   *Node
	oldest   *Node
	cache    map[string]*Node
	mu       sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	oldest := NewNode("", CacheEntry{})
	latest := NewNode("", CacheEntry{})
	oldest.next = latest
	latest.prev = oldest

	return &LRUCache{
		capacity: capacity,
		latest:   latest,
		oldest:   oldest,
		cache:    make(map[string]*Node),
	}
}

func (c *LRUCache) Get(url string) (CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, found := c.cache[url]; found {
		c.remove(node)
		c.add(node)
		return node.value, true
	}
	return CacheEntry{}, false
}

func (c *LRUCache) Put(url string, value CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, found := c.cache[url]; found {
		node.value = value
		c.remove(node)
		c.add(node)
		return
	}

	newNode := NewNode(url, value)
	c.cache[url] = newNode
	c.add(newNode)

	if len(c.cache) > c.capacity {
		lruNode := c.oldest.next
		delete(c.cache, lruNode.url)
		c.remove(lruNode)
	}
}

func (c *LRUCache) add(node *Node) {
	node.next = c.latest
	node.prev = c.latest.prev
	c.latest.prev.next = node
	c.latest.prev = node
}

func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
}
