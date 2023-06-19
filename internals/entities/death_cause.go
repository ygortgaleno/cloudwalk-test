package entities

import "sync"

type DeathCause struct {
	Name    string
	Counter uint
	sync    sync.RWMutex
}

func (dc *DeathCause) UpdateCounter() {
	dc.sync.Lock()
	defer dc.sync.Unlock()

	dc.Counter++
}
