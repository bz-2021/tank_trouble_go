package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"log"
	"math"
)

func main() {
	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Tank Trouble")
	if err := ebiten.RunGame(&Game{
		Mode:    ModeTitle,
		TankA:   NewTank(colornames.Grey, 140, 100, math.Pi),
		TankB:   NewTank(Red, 820, 520, 0),
		GameMap: NewMap(),
	}); err != nil {
		log.Fatal(err)
	}
}
