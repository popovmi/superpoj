package warsgame

import (
	"math"
	"reflect"
	"sync"
)

//go:generate msgp

type Game struct {
	Players     map[string]*Player `msg:"players"`
	CId         string             `msg:"cId"`
	PortalLinks []*PortalLink      `msg:"portalLinks"`
	Bricks      []*Brick           `msg:"bricks"`

	mu sync.Mutex
}

func NewGame() *Game {
	return &Game{
		Players: make(map[string]*Player),
		PortalLinks: []*PortalLink{
			NewPortalLink(350, 350, FieldWidth-350, FieldHeight-350),
			NewPortalLink(350, FieldHeight-350, FieldWidth-350, 350),
		},
		Bricks: []*Brick{
			{769, 679, 200, 40, "H"},
			{1031, 781, 200, 40, "H"},
			{769, 781, 200, 40, "H"},
			{1031, 679, 200, 40, "H"},
			{980, 217, 40, 200, "V"},
			{980, 417, 40, 200, "V"},
			{980, 883, 40, 200, "V"},
			{980, 1083, 40, 200, "V"},
		},
	}
}

func (g *Game) SetPlayers(players map[string]*Player) {
	g.Players = players
}

func (g *Game) AddPlayer(p *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.findFreeSpot(p)
	g.Players[p.ID] = p
	if len(g.Players) == 1 {
		g.setChaser(p)
	}
}

func (g *Game) findFreeSpot(np *Player) {
	if len(g.Players) == 0 {
		np.X = Left
		np.Y = Top
		return
	}

	for y := float64(Top); y <= Bottom; y += 1 {
		for x := float64(Left); x <= Right; x += 1 {
			np.X = x
			np.Y = y
			intersects := false
			for _, p := range g.Players {
				if np.Touching(p) {
					intersects = true
					break
				}
			}
			if !intersects {
				return
			}
		}
	}
}

func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	p := g.Players[id]
	delete(g.Players, id)
	if pV := reflect.ValueOf(p); !pV.IsNil() && p.ID == g.CId {
		if len(g.Players) == 0 {
			g.CId = ""
		} else {
			for _, p := range g.Players {
				g.setChaser(p)
				break
			}
		}
	}
}

func (g *Game) Tick() {
	g.mu.Lock()
	defer g.mu.Unlock()
	l := len(g.Players)
	if l > 0 {
		for _, p := range g.Players {
			p.Tick()
		}
	}
	g.detectCollisions()
}

func (g *Game) Teleport(id string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if p, ok := g.Players[id]; ok {
		for _, link := range g.PortalLinks {
			if link.CollideAndTeleport(p) {
				return true
			}
		}
	}
	return false
}

func (g *Game) detectCollisions() {
	for k1, p1 := range g.Players {
		for _, brick := range g.Bricks {
			brick.CollideAndBounce(p1)
		}
		for k2, p2 := range g.Players {
			if k1 < k2 {
				if p1.Touching(p2) {

					dx := p2.X - p1.X
					dy := p2.Y - p1.Y
					d := math.Sqrt(dx*dx + dy*dy)

					nx, ny := dx/d, dy/d
					tx, ty := -ny, nx

					v1n := p1.Vx*nx + p1.Vy*ny
					v1t := p1.Vx*tx + p1.Vy*ty
					v2n := p2.Vx*nx + p2.Vy*ny
					v2t := p2.Vx*tx + p2.Vy*ty

					v1n, v2n = v2n, v1n

					p1.Vx = (v1n*nx + v1t*tx) * PlayerElasticity
					p1.Vy = (v1n*ny + v1t*ty) * PlayerElasticity

					p2.Vx = (v2n*nx + v2t*tx) * PlayerElasticity
					p2.Vy = (v2n*ny + v2t*ty) * PlayerElasticity

					speedLimit(&p1.Vx, MaxCollideVelocity)
					speedLimit(&p1.Vy, MaxCollideVelocity)
					speedLimit(&p2.Vx, MaxCollideVelocity)
					speedLimit(&p2.Vy, MaxCollideVelocity)

					dd := float64(2*Radius) - d

					p1.X = p1.X - (nx)*(dd/2)
					p1.Y = p1.Y - (ny)*(dd/2)

					p2.X = p2.X + (nx)*(dd/2)
					p2.Y = p2.Y + (ny)*(dd/2)

					var newChaser *Player
					if p1.ID == g.CId {
						newChaser = p2
					} else if p2.ID == g.CId {
						newChaser = p1
					}
					if newChaser != nil {
						g.setChaser(newChaser)
					}
				}
			}
		}

	}
}

func (g *Game) setChaser(p *Player) {
	g.CId = p.ID
	p.ChaseCount++
}
