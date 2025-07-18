package utils

import (
	"sync"
)

type SmartMutex[K comparable] struct {
	singleMutex sync.Map
}

func NewSmartMutex[K comparable]() *SmartMutex[K] {
	return &SmartMutex[K]{
		singleMutex: sync.Map{},
	}
}

func (ctx *SmartMutex[K]) Lock(id K) {
	if val, ok := ctx.singleMutex.Load(id); ok {
		val.(*sync.Mutex).Lock()
	} else {
		newMutex := &sync.Mutex{}
		newMutex.Lock()
		ctx.singleMutex.Store(id, newMutex)
	}
}

func (ctx *SmartMutex[K]) Unlock(id K) {
	val, ok := ctx.singleMutex.Load(id)
	if !ok {
		return
	}
	val.(*sync.Mutex).Unlock()
}
