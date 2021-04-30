package lfu

import "container/list"

type cacheItem struct {
	key             string
	value           interface{}
	frequencyParent *list.Element
}

type frequencyItem struct {
	entries   map[*cacheItem]byte
	frequency int
}

type Cache struct {
	byKey       map[string]*cacheItem
	frequencies *list.List
	capacity    int
	size        int
}

func New(capacity int) *Cache {
	cache := new(Cache)
	cache.byKey = make(map[string]*cacheItem)
	cache.frequencies = list.New()
	cache.size = 0
	cache.capacity = capacity

	return cache
}

func (c *Cache) Set(key string, value interface{}) {
	if item, ok := c.byKey[key]; ok {
		item.value = value
		c.increment(item)

		return
	}

	item := new(cacheItem)
	item.key = key
	item.value = value
	c.byKey[key] = item
	c.size++

	if c.atCapacity() {
		c.evict(1)
	}

	c.increment(item)
}

func (c *Cache) Get(key string) interface{} {
	if item, ok := c.byKey[key]; ok {
		c.increment(item)
		return item
	}

	return nil
}

func (c *Cache) increment(item *cacheItem) {
	currentFreq := item.frequencyParent

	var nextFreqAmount int
	var nextFreq *list.Element

	if currentFreq == nil {
		nextFreqAmount = 1
		nextFreq = c.frequencies.Front()
	} else {
		nextFreqAmount = currentFreq.Value.(*frequencyItem).frequency + 1
		nextFreq = currentFreq.Next()
	}

	if nextFreq == nil || nextFreq.Value.(*frequencyItem).frequency != nextFreqAmount {
		newFreqItem := new(frequencyItem)
		newFreqItem.frequency = nextFreqAmount
		newFreqItem.entries = make(map[*cacheItem]byte)

		if currentFreq == nil {
			nextFreq = c.frequencies.PushFront(newFreqItem)
		} else {
			nextFreq = c.frequencies.InsertAfter(newFreqItem, currentFreq)
		}
	}

	item.frequencyParent = nextFreq
	nextFreq.Value.(*frequencyItem).entries[item] = 1
	if currentFreq != nil {
		c.remove(currentFreq, item)
	}
}

func (c *Cache) remove(listItem *list.Element, item *cacheItem) {
	frequencyItem := listItem.Value.(*frequencyItem)
	delete(frequencyItem.entries, item)

	if len(frequencyItem.entries) == 0 {
		c.frequencies.Remove(listItem)
	}
}

func (c *Cache) evict(count int) {
	for i := 0; i < count; {
		if item := c.frequencies.Front(); item != nil {
			for entry := range item.Value.(*frequencyItem).entries {
				if i < count {
					delete(c.byKey, entry.key)
					c.remove(item, entry)
					c.size--
					i++
				}
			}
		}
	}
}

func (c *Cache) atCapacity() bool {
	return c.size >= c.capacity
}
