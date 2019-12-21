# nearcache

nearcache is an in-memory key:value cache, this run in the golang environment using the same memory that program whitch 
are usage it.

Any kind of object can be stored into the cache, and can personalize the duration.

### Installation

`go get github.com/aeolabs/nearcache`

### Usage

```go
import (
  "fmt"
  "github.com/aeolabs/nearcache"
  "time"
)

func main(){
  ncache := InitNearCache()
  ncache.Add("key", "value", time.Second * 10)
  v, e := ncache.Get("key")    
  if e == nil {
     fmt.Printf("Cache value [%s]\n", v)
  } else{
  	 fmt.Println("error")
  }
}
```

### Future work

- [x] basic operations (Get, Add, Refresh, Del, Expire)
- [ ] add testing options
- [ ] evict functions
- [ ] notifications
- [ ] stats
- [ ] redis support (pub/sub)
- [ ] sync with others clients (using ae-net-prot 1.0)


