package nearcache

import (
	"errors"
	"sync"
	"time"
)

type NearCache struct {
	mux   sync.Mutex
	items map[string]Value
}

var (
	ErrNoExists = errors.New("no exists item")
	ErrExpire   = errors.New("item has been expired")
)

//InitNearCache a simple nearcache without configuration parameters
// ncache := nearcache.InitNearCache()
// ncache.Add("v1", "v1", time.Seconds * 5)
// item := ncache.Get("v1")
// fmt.println(item)
func InitNearCache() *NearCache {
	return &NearCache{
		items: make(map[string]Value),
	}
}

// Add a new item to the map, this value must be usage with duration
func (n *NearCache) Add(key string, value interface{}, duration time.Duration) error {
	n.mux.Lock()
	v := Value{
		value:    value,
		expire:   time.Now().Add(duration).UnixNano(),
		duration: duration,
	}
	n.items[key] = v
	n.mux.Unlock()
	return nil
}

//Get value from cache, if the item is expire or does not exists, this return an error
func (n *NearCache) Get(key string) (interface{}, error) {
	return n.get(key)
}

//Get Item and then expire
func (n *NearCache) GetAndExpire(key string) (interface{}, error) {
	v, e := n.get(key)
	if e == nil {
		n.expire(key)
	} else {
		return nil, ErrNoExists
	}
	return v, nil
}

//Get item and refresh the expiration time
func (n *NearCache) GetAndRefresh(key string) (interface{}, error) {
	return n.refresh(key)
}

func (n *NearCache) get(key string) (interface{}, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if &v == nil {
		return nil, ErrNoExists
	}
	if v.expire > time.Now().UnixNano() {
		return v.value, nil
	} else {
		delete(n.items, key)
		return nil, ErrExpire
	}

}

//Determine if the value in the cache is expired or not
// if the value is expired this return true, otherwise false
func (n *NearCache) Expired(key string) (bool, error) {
	return n.expire(key)
}

func (n *NearCache) expire(key string) (bool, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if &v == nil {
		return false, ErrNoExists
	}
	return v.expire > time.Now().UnixNano(), nil
}

//Delete the item from cache if its exists.
func (n *NearCache) Del(key string) error {
	return n.del(key)
}

func (n *NearCache) del(key string) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	_, e := n.get(key)
	if e != nil {
		return ErrNoExists
	}
	delete(n.items, key)
	return nil
}

//Refresh item into cache using configuration when this were added
func (n *NearCache) Refresh(key string) (interface{}, error) {
	return n.refresh(key)
}

func (n *NearCache) refresh(key string) (interface{}, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if &v == nil {
		return nil, ErrNoExists
	}
	v.expire = time.Now().Add(v.duration).UnixNano()
	return v.value, nil
}

//Update items into cache and return new value
func (n *NearCache) Update(key string, value interface{}) (interface{}, error) {
	return n.update(key, value)
}

func (n *NearCache) update(key string, value interface{}) (interface{}, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if &v == nil {
		return nil, ErrNoExists
	}
	v.value = value
	v.expire = time.Now().Add(v.duration).UnixNano()
	return value, nil
}
