[![GitHub release (latest by date)](https://img.shields.io/github/v/release/aeolabs/nearcache)](https://github.com/aeolabs/nearcache/releases)

# nearcache

nearcache is an in-memory key:value cache, this run in the golang environment using the same memory that program whitch 
are usage it.

Any kind of object can be stored into the cache, and can personalize the duration.

### Installation

`go get github.com/aeolabs/nearcache`

### Usage

```go
package main

import (
  "fmt"
  "github.com/aeolabs/nearcache"
  "time"
)

func main(){
  ncache := InitNearCache()
  ncache.OnDelete = func () (interface{}, error) {
     fmt.Println("element were deleted")
     return nil, nil
  }
  ncache.Add("key", "value", time.Second * 10)
  v, e := ncache.Get("key")    
  if e == nil {
     fmt.Printf("Cache value [%s]\n", v)
     ncache.Del("key") //this line should print "element were deleted"
  } else{
  	 fmt.Println("error")
  }
  
  items := ncache.Count()
  fmt.Println(items)
  
}
```

### Future work

- [x] basic operations (Get, Add, Refresh, Del, Expire)
- [x] count elements
- [ ] elements events (using bus events)
- [ ] add testing options
- [ ] evict functions
- [ ] notifications
- [ ] stats
- [ ] cache support (pub/sub)


