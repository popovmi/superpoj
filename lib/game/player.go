package warsgame

import (
	"log"
	"math"
	"strings"
	"sync"

	"github.com/matoous/go-nanoid/v2"

	"wars/lib/color"
)

//go:generate msgp

type Player struct {
	ID         string     `msg:"id" json:"id"`
	Name       string     `msg:"name" json:"name"`
	X          float64    `msg:"x" json:"x"`
	Y          float64    `msg:"y" json:"y"`
	Vx         float64    `msg:"vx" json:"vx"`
	Vy         float64    `msg:"vy" json:"vy"`
	Direction  string     `msg:"dir" json:"dir"`
	Color      color.RGBA `msg:"clr" json:"clr"`
	JoinedAt   int64      `msg:"joinedAt" json:"joinedAt"`
	ChaseCount int        `msg:"chaseCount" json:"chaseCount"`

	mu sync.Mutex
}

func NewPlayer() *Player {
	id, err := gonanoid.New()
	if err != nil {
		log.Fatal(err)
	}
	return &Player{ID: id}
}

func (p *Player) Tick() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Direction != "" {
		var dvx, dvy float64

		switch p.Direction {
		case "l":
			dvx -= Acceleration
		case "r":
			dvx += Acceleration
		case "u":
			dvy -= Acceleration
		case "d":
			dvy += Acceleration
		case "lu":
			dvx -= Acceleration / math.Sqrt2
			dvy -= Acceleration / math.Sqrt2
		case "ru":
			dvx += Acceleration / math.Sqrt2
			dvy -= Acceleration / math.Sqrt2
		case "ld":
			dvx -= Acceleration / math.Sqrt2
			dvy += Acceleration / math.Sqrt2
		case "rd":
			dvx += Acceleration / math.Sqrt2
			dvy += Acceleration / math.Sqrt2
		}

		accelerate(&p.Vx, dvx)
		accelerate(&p.Vy, dvy)
	}

	step(&p.X, &p.Vx)
	step(&p.Y, &p.Vy)

	wallReflect(&p.X, &p.Vx, Left, Right)
	wallReflect(&p.Y, &p.Vy, Top, Bottom)

	frictionBrake(&p.Vx)
	frictionBrake(&p.Vy)
}

func accelerate(v *float64, a float64) {
	if math.Abs(*v) < MaxVelocity {
		*v += a
		speedLimit(v, MaxVelocity)
	}
}

func speedLimit(v *float64, maxV float64) {
	if math.Abs(*v) > maxV {
		*v = math.Copysign(maxV, *v)
	}
}

func step(pos, v *float64) {
	*pos += *v
}

func wallReflect(pos, v *float64, min, max float64) {
	if *pos < min {
		*pos = min
		*v = -*v * WallElasticity
	} else if *pos > max {
		*pos = max
		*v = -*v * WallElasticity
	}

	speedLimit(v, MaxCollideVelocity)
}

func frictionBrake(v *float64) {
	if *v != 0 {
		if math.Abs(*v) > Friction {
			*v -= math.Copysign(Friction, *v)
		} else {
			*v = 0
		}
	}
}

func (p *Player) Move(dir string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Direction = dir
}

func (p *Player) Touching(p2 *Player) bool {
	d := math.Sqrt(math.Pow(p.X-p2.X, 2) + math.Pow(p.Y-p2.Y, 2))
	return d <= 2*Radius
}

func (p *Player) Brake() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Direction = ""
	p.Vx = p.Vx * Braking
	p.Vy = p.Vy * Braking
}

func (p *Player) Compare(ap *Player) int {
	return strings.Compare(p.ID, ap.ID)
}
