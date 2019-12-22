package nearcache

import "time"

type cacheitem struct {
	value    interface{}
	expire   int64
	duration time.Duration
}

type EventType int

const (
	OnDeleteEvt EventType = iota
	OnRefershEvt
	OnAddEvt
	OnUpdateEvt
	OnExpireEvt
)
