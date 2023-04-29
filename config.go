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

import "errors"

type Config struct {
	OnDelete  event
	OnUpdate  event
	OnAdd     event
	OnRefresh event
	OnExpire  event
}

type EventType int

const (
	OnDeleteEvt EventType = iota
	OnRefershEvt
	OnAddEvt
	OnUpdateEvt
	OnExpireEvt
)

func (cf *Config) doCommand(evt EventType, item *Cacheitem) (interface{}, error) {
	switch evt {
	case OnAddEvt:
		if cf.OnAdd != nil {
			return cf.OnAdd(item)
		}
	case OnDeleteEvt:
		if cf.OnDelete != nil {
			return cf.OnDelete(item)
		}
	case OnRefershEvt:
		if cf.OnRefresh != nil {
			return cf.OnDelete(item)
		}
	case OnUpdateEvt:
		if cf.OnUpdate != nil {
			return cf.OnRefresh(item)
		}
	case OnExpireEvt:
		if cf.OnExpire != nil {
			return cf.OnExpire(item)
		}
	}
	return nil, errors.New("no events where defined")
}
