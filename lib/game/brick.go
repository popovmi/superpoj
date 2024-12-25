package warsgame

import "math"

//go:generate msgp

type Brick struct {
	X float64 `msg:"X"`
	Y float64 `msg:"Y"`
	W float64 `msg:"W"`
	H float64 `msg:"H"`
	D string  `msg:"D"`
}

func (b *Brick) CollideAndBounce(plr *Player) {
	nextX := plr.X + plr.Vx
	nextY := plr.Y + plr.Vy

	closestX := math.Max(b.X, math.Min(plr.X, b.X+b.W))
	closestY := math.Max(b.Y, math.Min(plr.Y, b.Y+b.H))

	distance := math.Sqrt(math.Pow(nextX-closestX, 2) + math.Pow(nextY-closestY, 2))

	if distance > Radius {
		return
	}

	nx := nextX - closestX
	ny := nextY - closestY

	if math.Abs(nx) > math.Abs(ny) {
		plr.Vx *= -1 * BrickElasticity
		plr.X = closestX + (Radius+0.0001)*math.Copysign(1, nx)
	} else if math.Abs(nx) < math.Abs(ny) {
		plr.Vy *= -1 * BrickElasticity
		plr.Y = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	} else {
		plr.Vx *= -1 * BrickElasticity
		plr.Vy *= -1 * BrickElasticity
		plr.X = closestX + (Radius+0.0001)*math.Copysign(1, nx)
		plr.Y = closestY + (Radius+0.0001)*math.Copysign(1, ny)
	}
}
