package ttlcache

import "time"

func (t *TTLCache) Get(key string) *KV {
	if k, ok := t.data[key]; ok {
		return &KV{
			Key:   key,
			Value: k.Value,
			TTL:   t.ttl - (time.Since(k.updated)),
		}
	}

	return nil
}

func (t *TTLCache) GetAll() []*KV {
	keys := []*KV{}
	for k, v := range t.data {
		keys = append(keys, &KV{
			Key:   k,
			Value: v,
			TTL:   t.ttl - (time.Since(v.updated)),
		})
	}
	return keys
}
