package nearcache

import (
	"errors"
	"sync"
	"time"
)

type event = func() (interface{}, error)

type NearCache struct {
	mux       sync.Mutex
	items     map[string]*cacheitem
	OnDelete  event
	OnAdd     event
	OnRefresh event
	OnUpdate  event
	OnExpire  event
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
		items: make(map[string]*cacheitem),
	}
}

// Add a new item to the map, this value must be usage with duration
func (n *NearCache) Add(key string, value interface{}, duration time.Duration) error {
	n.mux.Lock()
	v := &cacheitem{
		value:    value,
		expire:   time.Now().Add(duration).UnixNano(),
		duration: duration,
	}
	n.items[key] = v
	n.doCommand(OnAddEvt)
	n.mux.Unlock()
	return nil
}

//Get value from cache, if the item is expire or does not exists, this return an error
func (n *NearCache) Get(key string) (interface{}, error) {
	get, err := n.get(key)
	if err != nil {
		return nil, err
	}
	return get.value, nil
}

//Get Item and then expire
func (n *NearCache) GetAndExpire(key string) (interface{}, error) {
	v, e := n.get(key)
	if e == nil {
		n.expire(key)
	} else {
		return nil, ErrNoExists
	}
	return v.value, nil
}

//Get item and refresh the expiration time
func (n *NearCache) GetAndRefresh(key string) (interface{}, error) {
	return n.refresh(key)
}

func (n *NearCache) get(key string) (*cacheitem, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if v == nil {
		return nil, ErrNoExists
	}
	if v.expire > time.Now().UnixNano() {
		return v, nil
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
	if v == nil {
		return false, ErrNoExists
	}
	expired := v.expire > time.Now().UnixNano()
	if expired {
		n.doCommand(OnExpireEvt)
	}
	return expired, nil
}

//Delete the item from cache if its exists.
func (n *NearCache) Del(key string) error {
	return n.del(key)
}

func (n *NearCache) del(key string) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if v == nil {
		return ErrNoExists
	}
	n.doCommand(OnDeleteEvt)
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
	if v == nil {
		return nil, ErrNoExists
	}
	v.expire = time.Now().Add(v.duration).UnixNano()
	n.doCommand(OnRefershEvt)
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
	if v == nil {
		return nil, ErrNoExists
	}
	v.value = value
	v.expire = time.Now().Add(v.duration).UnixNano()
	n.doCommand(OnUpdateEvt)
	return value, nil
}

// Clean all the items into cache
func (n *NearCache) Clean() {
	n.mux.Lock()
	n.items = make(map[string]*cacheitem)
	n.mux.Unlock()
}

func (n *NearCache) doCommand(evt EventType) error {
	switch evt {
	case OnAddEvt:
		if n.OnAdd != nil {
			n.OnAdd()
		}
		break
	case OnDeleteEvt:
		if n.OnDelete != nil {
			n.OnDelete()
		}
		break
	case OnRefershEvt:
		if n.OnRefresh != nil {
			n.OnDelete()
		}
		break
	case OnUpdateEvt:
		if n.OnUpdate != nil {
			n.OnRefresh()
		}
		break
	case OnExpireEvt:
		if n.OnExpire != nil {
			n.OnExpire()
		}
		break
	}
	return nil
}
