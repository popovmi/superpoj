package warsgame

import (
	"math"
	"time"
)

//go:generate msgp

type Portal struct {
	X float64 `msg:"x"`
	Y float64 `msg:"y"`
}

type PortalLink struct {
	P1       *Portal          `msg:"p1"`
	P2       *Portal          `msg:"p2"`
	LastUsed map[string]int64 `msg:"LastUsed"`
}

func NewPortalLink(x1, y1, x2, y2 float64) *PortalLink {
	return &PortalLink{&Portal{x1, y1}, &Portal{x2, y2}, make(map[string]int64)}
}

func (p *Portal) Touching(plr *Player) bool {
	d := math.Sqrt(math.Pow(p.X-plr.X, 2) + math.Pow(p.Y-plr.Y, 2))
	return d <= (PortalRadius - Radius)
}

func (p *PortalLink) CollideAndTeleport(plr *Player) bool {
	if lu, used := p.LastUsed[plr.ID]; used {
		if lu > time.Now().Unix()-5 {
			return false
		}
	}
	ported := p.P1.CollideAndTeleport(plr, p.P2)
	if !ported {
		ported = p.P2.CollideAndTeleport(plr, p.P1)
	}
	if ported {
		p.LastUsed[plr.ID] = time.Now().Unix()
	}
	return ported
}

func (p *Portal) CollideAndTeleport(plr *Player, dest *Portal) bool {
	if !p.Touching(plr) {
		return false
	}

	plr.X = dest.X
	plr.Y = dest.Y

	return true
}
