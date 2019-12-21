package nearcache

import "time"

type Value struct {
	value    interface{}
	expire   int64
	duration time.Duration
}
