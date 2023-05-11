package nearcache

//	MIT License
//
//	Copyright (c) Aeolabs SRL. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// If it is usage in cache item, at this point is possible to know each value
type event = func(item *Cacheitem) (interface{}, error)

type Stats struct {
	calls  uint64
	misses uint64
	added  uint64
	items  uint64
}

type NearCache struct {
	mux    sync.Mutex
	items  map[string]*Cacheitem
	config *Config
	stats  *Stats
}

var (
	ErrNoExists = errors.New("no exists item")
	ErrExpire   = errors.New("item has been expired")
)

// Init InitNearCache a simple nearcache without configuration parameters
// ncache := nearcache.InitNearCache()
// ncache.Add("v1", "v1", time.Seconds * 5)
// item := ncache.Get("v1")
// fmt.println(item)
func Init() *NearCache {
	cfg := &Config{}

	return &NearCache{
		items:  make(map[string]*Cacheitem),
		config: cfg,
		stats: &Stats{
			calls:  0,
			misses: 0,
			added:  0,
			items:  0,
		},
	}
}

func InitWithConfig(cfg *Config) *NearCache {
	return &NearCache{
		items:  make(map[string]*Cacheitem),
		config: cfg,
		stats: &Stats{
			calls:  0,
			misses: 0,
			added:  0,
			items:  0,
		},
	}
}

// Add a new item to the map, this value must be usage with duration
func (n *NearCache) Add(key string, value interface{}, duration time.Duration) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := &Cacheitem{
		Value:    value,
		expire:   time.Now().Add(duration).UnixNano(),
		duration: duration,
	}
	n.items[key] = v
	n.config.doCommand(OnAddEvt, v)
	atomic.AddUint64(&n.stats.added, 1)
	return nil
}

// Get value from cache, if the item is expire or does not exists, this return an error
func (n *NearCache) Get(key string) (*Cacheitem, error) {
	citem, err := n.get(key)
	if err != nil {
		return nil, err
	}
	return citem, nil
}

// GetAndExpire Get Item and then expire
func (n *NearCache) GetAndExpire(key string) (*Cacheitem, error) {
	v, e := n.get(key)
	if e == nil {
		n.expire(key)
	} else {
		return nil, ErrNoExists
	}
	return v, nil
}

// GetAndRefresh Get item and refresh the expiration time
func (n *NearCache) GetAndRefresh(key string) (*Cacheitem, error) {
	return n.refresh(key)
}

func (n *NearCache) get(key string) (*Cacheitem, error) {
	v := n.items[key]
	if v == nil {
		atomic.AddUint64(&n.stats.misses, 1)
		return nil, ErrNoExists
	}
	if v.expire > time.Now().UnixNano() {
		atomic.AddUint64(&n.stats.calls, 1)
		return v, nil
	} else {
		n.cleanItem(key)
		atomic.AddUint64(&n.stats.misses, 1)
		return nil, ErrExpire
	}
}

func (n *NearCache) Has(key string) bool {
	_, ok := n.items[key]
	return ok
}

// Expired Determine if the value in the cache is expired or not
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
	return v.Expired()
}

// Del Delete the item from cache if its exists.
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
	n.config.doCommand(OnDeleteEvt, v)
	n.cleanItem(key)
	return nil
}

// Refresh item into cache using configuration when this were added
func (n *NearCache) Refresh(key string) (*Cacheitem, error) {
	return n.refresh(key)
}

func (n *NearCache) refresh(key string) (*Cacheitem, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if v == nil {
		return nil, ErrNoExists
	}
	v.refersh()
	return v, nil
}

// Update items into cache and return new value
func (n *NearCache) Update(key string, value interface{}) (*Cacheitem, error) {
	return n.update(key, value)
}

func (n *NearCache) update(key string, value interface{}) (*Cacheitem, error) {
	n.mux.Lock()
	defer n.mux.Unlock()
	v := n.items[key]
	if v == nil {
		return nil, ErrNoExists
	}
	return v.update(v), nil
}

func (n *NearCache) cleanItem(key string) {
	delete(n.items, key)
	atomic.AddUint64(&n.stats.items, -1)
}

// Clean all the items into cache
func (n *NearCache) Clean() {
	n.mux.Lock()
	defer n.mux.Unlock()
	n.items = make(map[string]*Cacheitem)
	n.stats.items = 0
}

func (n *NearCache) Count() uint64 {
	n.mux.Lock()
	defer n.mux.Unlock()
	return n.stats.items
}

// Statics Get stats about cache
func (n *NearCache) Statics() *Stats {
	return n.stats
}
