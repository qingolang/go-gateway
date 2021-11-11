package common

import (
	"sync"
	"time"
)

// FlowCounterHandler
var FlowCounterHandler *FlowCounter

// FlowCounter
type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

// NewFlowCounter
func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

// init
func init() {
	FlowCounterHandler = NewFlowCounter()
}

// GetCounter
func (counter *FlowCounter) GetCounter(serverName string) (*RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serverName {
			return item, nil
		}
	}

	newCounter := NewRedisFlowCountService(serverName, 1*time.Second)
	counter.Locker.Lock()
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.RedisFlowCountMap[serverName] = newCounter
	counter.Locker.Unlock()
	return newCounter, nil
}
