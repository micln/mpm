package cache

import "time"

type cacheData struct {
	Value    interface{}
	FailedAt time.Time
}

var data map[string]*cacheData

func init() {
	Clear()
}

func Get(key string) interface{} {
	return data[key].Value
}

func Set(key string, val interface{}, ttls ...int) {
	ttl := 300
	if len(ttls) > 0 {
		ttl = ttls[0]
	}

	data[key] = &cacheData{
		val,
		time.Now().Add(ttl * time.Second),
	}

	go autoDel(key, ttl)

}

func autoDel(key string, ttl int) {
	time.Sleep(ttl * time.Second)
	Del(key)
}

func Del(key string) {
	data[key] = nil
}

func Clear() {
	data = make(map[string]interface{})
}
