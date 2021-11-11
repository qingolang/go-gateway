package common

import (
	"sync"

	"golang.org/x/time/rate"
)

var FlowLimiterHandler *FlowLimiter

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

// FlowLimiter
type FlowLimiter struct {
	FlowLmiterMap   map[string]*FlowLimiterItem
	FlowLmiterSlice []*FlowLimiterItem
	Locker          sync.RWMutex
}

// FlowLimiterItem
type FlowLimiterItem struct {
	ServiceName string
	Limter      *rate.Limiter
}

// NewFlowLimiter
func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLmiterMap:   map[string]*FlowLimiterItem{},
		FlowLmiterSlice: []*FlowLimiterItem{},
		Locker:          sync.RWMutex{},
	}
}

// GetLimiter
func (counter *FlowLimiter) GetLimiter(serverName string, qps float64) (*rate.Limiter, error) {
	for _, item := range counter.FlowLmiterSlice {
		if item.ServiceName == serverName {
			return item.Limter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		ServiceName: serverName,
		Limter:      newLimiter,
	}
	counter.Locker.Lock()
	counter.FlowLmiterSlice = append(counter.FlowLmiterSlice, item)
	counter.FlowLmiterMap[serverName] = item
	counter.Locker.Unlock()
	return newLimiter, nil
}
