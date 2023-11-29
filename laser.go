package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

const LaserLength = 300
const OneStep = 1

type Laser struct {
	towards float64
	valid   bool
}

func (l *Laser) DrawLaser(screen *ebiten.Image, tank *Tank) {
	if !l.valid {
		return
	}
	l.towards = math.Pi + tank.Angle
	x, y := tank.x, tank.y
	for i := 0; i < LaserLength; i++ {
		if isBlack(screen.At(int(x), int(y))) {
			if isBlack(screen.At(int(x-5), int(y))) && isBlack(screen.At(int(x+5), int(y))) {
				l.towards = -l.towards
			} else if isBlack(screen.At(int(x), int(y-5))) && isBlack(screen.At(int(x), int(y+5))) {
				l.towards = math.Pi - l.towards
			}
		}
		x, y = x+OneStep*math.Cos(l.towards), y+OneStep*math.Sin(l.towards)

		//ran := rand.Int() % 100
		//if ran > 10 {
		screen.Set(int(x-1), int(y-1), Red)
		//}
	}
}

func isBlack(color color.Color) bool {
	r, g, b, _ := color.RGBA()
	if r < 0x80 && g < 0x80 && b < 0x80 {
		return true
	}
	return false
}
