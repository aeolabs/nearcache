package nearcache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleAddCache(t *testing.T) {
	ncache := InitNearCache()
	ncache.OnDelete = func() (i interface{}, err error) {
		fmt.Println("item deleted")
		return nil, nil
	}

	key := "test1"
	ncache.Add(key, "test", time.Second*60)
	v, e := ncache.Get("test1")
	ncache.Del(key)
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
	fmt.Println(e)
}

func TestAddAndRefresh(t *testing.T) {
	ncache := InitNearCache()
	key := "test1"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test")
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
	v, _ = ncache.refresh(key)
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
}

func TestAddAndModify(t *testing.T) {
	ncache := InitNearCache()
	key := "test"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test1")
	assert.Equal(t, v, "test")
	update, _ := ncache.Update(key, "test2")
	assert.Equal(t, update, "test2")
}

func TestCleanItems(t *testing.T) {
	ncache := InitNearCache()
	key := "test"
	ncache.Add(key, "value", time.Second*10)
	ncache.Get(key)
	ncache.Clean()
	get, err := ncache.Get(key)
	assert.Equal(t, get, nil)
	assert.NotNil(t, err)
}
