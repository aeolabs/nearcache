package nearcache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleAddCache(t *testing.T) {
	ncache := InitNearCache()
	key := "test1"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test1")
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
}

func TestAddAndRefresh(t *testing.T) {
	ncache := InitNearCache()
	key := "test1"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test1")
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
	v, _ = ncache.refresh(key)
	fmt.Printf(" cache key[%s] with value [%s]\n", key, v)
}

func TestAddAndModify(t *testing.T) {
	ncache := InitNearCache()
	key := "test1"
	ncache.Add(key, "test", time.Second*10)
	v, _ := ncache.Get("test1")
	assert.Equal(t, v, "test")
	update, _ := ncache.Update(key, "test2")
	assert.Equal(t, update, "test2")
}
