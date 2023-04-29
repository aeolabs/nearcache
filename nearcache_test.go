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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleAddCache(t *testing.T) {
	cfg := &Config{
		OnDelete: func(n *Cacheitem) (i interface{}, err error) {
			fmt.Printf("item [%s] were deleted\n", n.Value)
			return nil, nil
		},
	}
	ncache := InitWithConfig(cfg)

	key := "test1"
	ncache.Add(key, "test", time.Second*60)
	v, e := ncache.Get("test1")
	ncache.Del(key)
	fmt.Printf(" cache key[%s] with value [%v]\n", key, v)
	fmt.Println(e)
}

func TestAddAndRefresh(t *testing.T) {
	ncache := Init()
	key := "test1"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test")
	fmt.Printf(" cache key[%s] with value [%v]\n", key, v)
	v, _ = ncache.refresh(key)
	fmt.Printf(" cache key[%s] with value [%v]\n", key, v)
}

func TestAddAndModify(t *testing.T) {
	ncache := Init()
	key := "test"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get(key)
	assert.Equal(t, v.Value, "test")
	update, _ := ncache.Update(key, "test2")
	assert.Equal(t, update.Value, "test2")
}

func TestCleanItems(t *testing.T) {
	ncache := Init()
	key := "test"
	ncache.Add(key, "value", time.Second*10)
	ncache.Get(key)
	ncache.Clean()
	get, err := ncache.Get(key)
	assert.Equal(t, get, nil)
	assert.NotNil(t, err)
}

func TestHasItem(t *testing.T) {
	ncache := Init()
	key := "test"
	ncache.Add(key, "value", time.Second*10)
	has := ncache.Has(key)
	assert.Equal(t, has, true, "the elements is not in the cache")
}

func TestCountItems(t *testing.T) {
	ncache := Init()
	key := "test"
	ncache.Add(key, "value", time.Second*10)
	count := ncache.Count()
	assert.Equal(t, 1, count, "No elements")
}
