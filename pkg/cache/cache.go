package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var mem = memcache.New(":11211")

func Get(key string) (*memcache.Item, error) {
	item, err := mem.Get(key)

	// item dont setted
	if err != nil && item == nil {
		return nil, nil
		// error
	} else if err != nil {
		return nil, err
		// item exist
	} else {
		return item, nil
	}
}

func Set(key string, value []byte, expiration int32) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
	return mem.Set(item)
}
