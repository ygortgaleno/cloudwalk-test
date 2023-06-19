package entities

import (
	"sync"
)

type Player struct {
	Name        string
	Kills       uint32
	WorldDeaths int64
	sync        sync.RWMutex
}

func (p *Player) UpdatePlayerKill() {
	p.sync.Lock()
	defer p.sync.Unlock()

	p.Kills += 1
}

func (p *Player) UpdateWorldDeaths(val int64) {
	p.sync.Lock()
	defer p.sync.Unlock()

	p.WorldDeaths += val
}
