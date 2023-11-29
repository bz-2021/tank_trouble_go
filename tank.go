package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

const (
	CircleRadius = 12
	TankHeight   = 42
	TankWidth    = 34
)

type Tank struct {
	x     float64
	y     float64
	Angle float64
	Vx    float64
	Vy    float64

	Color  color.Color
	Body   *ebiten.Image
	barrel *ebiten.Image

	Bullets     [5]*Bullet
	validBullet uint8

	TankAHorizontalCollision bool
	TankAVerticalCollision   bool

	laser Laser
}

func (t *Tank) GetX() float64 {
	return t.x
}

func (t *Tank) GetY() float64 {
	return t.y
}

func NewTank(col color.Color, xPos, yPos float64, tankAngle float64) *Tank {
	body := ebiten.NewImage(TankHeight, TankWidth)
	r, g, b, _ := col.RGBA()
	barrel := ebiten.NewImage(12, 36)
	barrel.Fill(col)
	bodyColor := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 0xcc,
	}
	body.Fill(bodyColor)

	var bullets [5]*Bullet
	for i, _ := range bullets {
		bulletBody := ebiten.NewImage(8, 8)
		for y := 0.0; y <= 8; y += 1 {
			l := math.Sqrt(16 - (4-y)*(4-y))
			for x := 4 - l; x <= 4+l; x += 1 {
				bulletBody.Set(int(x), int(y), color.Black)
			}
		}
		bullets[i] = &Bullet{
			X:         100,
			Y:         100,
			Vx:        0,
			Vy:        0,
			Valid:     false,
			lifeCycle: 0,
			body:      bulletBody,
		}
	}

	return &Tank{
		x:       xPos,
		y:       yPos,
		Color:   col,
		Body:    body,
		barrel:  barrel,
		Angle:   tankAngle,
		Bullets: bullets,
		laser: Laser{
			valid: true,
		},
	}
}

func (t *Tank) DrawTank(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	barrelOp := &ebiten.DrawImageOptions{}

	t.x += t.Vx * math.Cos(t.Angle)
	t.y += t.Vy * math.Sin(t.Angle)
	w, h := t.Body.Bounds().Dx(), t.Body.Bounds().Dy()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Rotate(t.Angle)
	op.GeoM.Translate(t.x, t.y)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(t.Body, op)

	barrelOp.GeoM.Translate(-float64(t.barrel.Bounds().Dx())/2.0, -float64(h)/2.0)
	barrelOp.GeoM.Rotate(math.Pi/2 + t.Angle)
	barrelOp.GeoM.Translate(t.x-8*math.Cos(t.Angle), t.y-8*math.Sin(t.Angle))
	screen.DrawImage(t.barrel, barrelOp)

	for y := t.y - CircleRadius; y <= CircleRadius+t.y; y += 1 {
		l := math.Sqrt(144 - (t.y-y)*(t.y-y))
		for x := t.x - l; x <= t.x+l; x += 1 {
			screen.Set(int(x), int(y), t.Color)
		}
	}
}

func (t *Tank) UpdateSpeed() {
	t.y += t.Vy
	t.x += t.Vx
	t.Vx, t.Vy = 0, 0
}

func (t *Tank) Fire(bullet *Bullet) {
	if bullet.Valid {
		return
	}
	bullet.Valid = true
	bullet.lifeCycle = BulletLifeCircle
	bullet.X = t.x - 15*math.Cos(t.Angle)
	bullet.Y = t.y - 15*math.Sin(t.Angle)
	bullet.Vx = -BulletLineSpeed * math.Cos(t.Angle)
	bullet.Vy = -BulletLineSpeed * math.Sin(t.Angle)

}
