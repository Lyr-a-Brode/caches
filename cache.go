package caches

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}
