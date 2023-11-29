package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BulletLineSpeed  = 4
	BulletLifeCircle = 270
)

type Bullet struct {
	X         float64
	Y         float64
	Vx        float64
	Vy        float64
	Valid     bool
	lifeCycle int32
	JustShoot int

	body      *ebiten.Image
	HitCount  int
	MissCount int
}

func (bullet *Bullet) UpdateBulletPos() {
	if bullet.lifeCycle < 0 {
		bullet.Valid = false
		bullet.JustShoot = 0
		return
	}
	bullet.lifeCycle--
	bullet.JustShoot++
	bullet.X += bullet.Vx
	bullet.Y += bullet.Vy
}

func (bullet *Bullet) DrawBullet(screen *ebiten.Image) {
	if bullet.Valid == false {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bullet.X-4, bullet.Y-4)
	screen.DrawImage(bullet.body, op)
}
