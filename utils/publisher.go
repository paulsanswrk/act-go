package utils

import (
	"golang.org/x/exp/maps"
	"sync"
)

type Publisher[T any] struct {
	subscribers map[uint64]func(T)
	nLastSub    uint64
}

func (pub *Publisher[T]) Add(sub func(T)) (nSub uint64) {
	var m sync.Mutex
	m.Lock()

	if pub.subscribers == nil {
		pub.subscribers = make(map[uint64]func(T))
	}

	nSub = pub.nLastSub
	pub.subscribers[pub.nLastSub] = sub
	pub.nLastSub++

	m.Unlock()

	return
}

func (pub *Publisher[T]) Remove(nSub uint64) {
	delete(pub.subscribers, nSub)
}

func (pub *Publisher[T]) RunAll(msg T) {
	for _, k := range maps.Keys(pub.subscribers) {
		go pub.subscribers[k](msg)
	}
}

func (pub *Publisher[T]) Any() bool {
	return len(pub.subscribers) > 0
}
