package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lfritz/mazes/generate"
	"github.com/lfritz/mazes/grid"
	"image/color"
	"math"
	"math/rand"
)

var (
	BackgroundColor = color.RGBA{R: 0xe6, G: 0xe6, B: 0xe6, A: 0xff}
	Black           = color.Black
)

const (
	MapCol        = 9
	MapRow        = 6
	CellWidth     = 760 / MapCol
	GlobalXOffset = 100
	GlobalYOffset = 60
)

type Map struct {
	Maze *grid.Maze
}

func NewMap() *Map {
	generated := grid.NewMaze(MapCol, MapRow, true)
	*generated.WallAbove(0, 0) = false
	*generated.WallAbove(generated.Width()-1, generated.Height()) = false
	generate.Backtracking(generated, rand.New(rand.NewSource(0)))

	return &Map{Maze: generated}
}

func (m *Map) DrawMap(screen *ebiten.Image) {

	//for x := 0; x < 200; x++ {
	//	for y := 0; y < 60; y++ {
	//		screen.Set(x, y, Black)
	//	}
	//}

	horizontalBlock := ebiten.NewImage(CellWidth+1, 3)
	horizontalBlock.Fill(Black)
	verticalBlock := ebiten.NewImage(3, CellWidth+1)
	verticalBlock.Fill(Black)
	op := &ebiten.DrawImageOptions{}

	var lineX = (m.Maze.Width() - 1) * CellWidth
	var lineY = m.Maze.Height() * CellWidth

	op.GeoM.Translate(GlobalXOffset, GlobalYOffset-1)
	screen.DrawImage(horizontalBlock, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(lineX+GlobalXOffset), float64(lineY+GlobalYOffset-1))
	screen.DrawImage(horizontalBlock, op)

	// draw horizontal walls
	for y := 0; y <= m.Maze.Height(); y++ {
		j := y * CellWidth
		for x := 0; x < m.Maze.Width(); x++ {
			if *m.Maze.WallAbove(x, y) {
				i := x * CellWidth
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i+GlobalXOffset), float64(j+GlobalYOffset-1))
				screen.DrawImage(horizontalBlock, op)
			}
		}
	}

	// draw vertical walls
	for x := 0; x <= m.Maze.Width(); x++ {
		i := x * CellWidth
		for y := 0; y < m.Maze.Height(); y++ {
			if *m.Maze.WallLeftOf(x, y) {
				j := y * CellWidth
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i+GlobalXOffset-1), float64(j+GlobalYOffset))
				screen.DrawImage(verticalBlock, op)
			}
		}
	}
}

func (m *Map) Hit(bullet *Bullet, x, y int) bool {

	var lineX = (m.Maze.Width() - 1) * CellWidth
	var lineY = m.Maze.Height() * CellWidth

	horizontalSign := false
	verticalSign := false

	for k := 0; k <= CellWidth; k++ {

		if math.Abs(bullet.X-float64(k+GlobalXOffset)) < 4 && math.Abs(bullet.Y-float64(GlobalYOffset)) < 4 {
			if bullet.Y > GlobalYOffset && bullet.Vy > 0 {
				horizontalSign = true
				break
			}
			bullet.Vy = -bullet.Vy
			break
		}

		if math.Abs(bullet.X-float64(lineX+k+GlobalXOffset)) < 4 && math.Abs(bullet.Y-float64(lineY+GlobalYOffset)) < 4 {
			if bullet.Y < float64(lineY+GlobalYOffset) && bullet.Vy < 0 {
				horizontalSign = true
				break
			}
			bullet.Vy = -bullet.Vy
			break
			//return true
		}
	}

	// draw horizontal walls
	for y := 0; y <= m.Maze.Height(); y++ {
		j := y * CellWidth
		for x := 0; x < m.Maze.Width(); x++ {
			if *m.Maze.WallAbove(x, y) {
				i := x * CellWidth
				for k := 0; k <= CellWidth; k++ {

					if math.Abs(bullet.X-float64(i+k+GlobalXOffset)) < 4 && math.Abs(bullet.Y-float64(j+GlobalYOffset)) < 4 {
						if bullet.Y > float64(j+GlobalYOffset) && bullet.Vy > 0 {
							horizontalSign = true
							break
						}
						if bullet.Y < float64(j+GlobalYOffset) && bullet.Vy < 0 {
							horizontalSign = true
							break
						}
						bullet.Vy = -bullet.Vy
						break
						//return true
					}
				}
			}
			if horizontalSign {
				break
			}
		}
		if horizontalSign {
			break
		}
	}

	// draw vertical walls
	for x := 0; x <= m.Maze.Width(); x++ {
		i := x * CellWidth
		for y := 0; y < m.Maze.Height(); y++ {
			if *m.Maze.WallLeftOf(x, y) {
				j := y * CellWidth
				for k := 0; k <= CellWidth; k++ {

					if math.Abs(bullet.X-float64(i+GlobalXOffset)) < 4 && math.Abs(bullet.Y-float64(j+k+GlobalYOffset)) < 4 {
						if bullet.X > float64(i+GlobalXOffset) && bullet.Vx > 0 {
							verticalSign = true
							break
						}
						if bullet.X < float64(i+GlobalXOffset) && bullet.Vx < 0 {
							verticalSign = true
							break
						}
						bullet.Vx = -bullet.Vx
						break
						//return true
					}
				}
			}
			if verticalSign {
				break
			}
		}
		if verticalSign {
			break
		}
	}
	return false
}

var cnt = 0

func (m *Map) CheckCollision(t *Tank) {

	HorizontalCollision, VerticalCollision := false, false

	tankCenterX, tankCenterY := t.GetX(), t.GetY()

	axis1First, axis1Second := math.Cos(t.Angle), -math.Sin(t.Angle)
	axis2First, axis2Second := math.Sin(t.Angle), math.Cos(t.Angle)

	var lineX = float64((m.Maze.Width() - 1) * CellWidth)
	var lineY = float64(m.Maze.Height() * CellWidth)

	centerX, centerY := GlobalXOffset+CellWidth/2-tankCenterX, GlobalYOffset-tankCenterY

	var Judge, Judge2 func() bool

	Judge = func() bool {
		projV := math.Abs(centerX*axis1First + centerY*axis1Second)
		proRadius := math.Abs(axis1First)*CellWidth/2 + math.Abs(axis1Second)*2
		if proRadius+TankHeight/2 <= projV {
			return false
		}
		projV = math.Abs(centerX*axis2First + centerY*axis2Second)
		proRadius = math.Abs(axis2First)*CellWidth/2 + math.Abs(axis2Second)*2
		if proRadius+TankWidth/2 <= projV {
			return false
		}
		projV = math.Abs(centerX)
		proRadius = math.Abs(axis1First)*TankHeight/2 + math.Abs(axis2First)*TankWidth/2
		if proRadius+CellWidth/2 <= projV {
			return false
		}
		projV = math.Abs(centerY)
		proRadius = math.Abs(axis1Second)*TankHeight/2 + math.Abs(axis2Second)*TankWidth/2
		if proRadius+2 <= projV {
			return false
		}
		if centerY <= 0 {
			if t.Vy < 0 {
				t.Vy = 0.5
			}
		} else {
			if t.Vy > 0 {
				t.Vy = -0.5
			}
		}
		if centerX+CellWidth/2 <= 0 && t.Vx < 0 {
			t.Vx = 0.5
		}
		if centerX-CellWidth/2 >= 0 && t.Vx > 0 {
			t.Vx = -0.5
		}
		return true
	}

	Judge2 = func() bool {
		projV := math.Abs(centerX*axis1First + centerY*axis1Second)
		proRadius := math.Abs(axis1Second)*CellWidth/2 + math.Abs(axis1First)
		if proRadius+TankHeight/2 <= projV {
			return false
		}
		projV = math.Abs(centerX*axis2First + centerY*axis2Second)
		proRadius = math.Abs(axis2Second)*CellWidth/2 + math.Abs(axis2First)
		if proRadius+TankWidth/2 <= projV {
			return false
		}
		projV = math.Abs(centerY)
		proRadius = math.Abs(axis1Second)*TankHeight/2 + math.Abs(axis2Second)*TankWidth/2
		if proRadius+CellWidth/2 <= projV {
			return false
		}
		projV = math.Abs(centerX)
		proRadius = math.Abs(axis1First)*TankHeight/2 + math.Abs(axis2First)*TankWidth/2
		if proRadius+1 <= projV {
			return false
		}
		if centerX <= 0 {
			if t.Vx < 0 {
				t.Vx = 0.5
			}
		} else {
			if t.Vx > 0 {
				t.Vx = -0.5
			}
		}
		if centerY+CellWidth/2 <= 0 && t.Vy < 0 {
			t.Vy = 0.5
		}
		if centerY-CellWidth/2 >= 0 && t.Vy > 0 {
			t.Vy = -0.5
		}
		return true
	}

	if Judge() {
		HorizontalCollision = true
	}

	centerX, centerY = lineX+GlobalXOffset+CellWidth/2-tankCenterX, lineY+GlobalYOffset-tankCenterY

	if Judge() {
		HorizontalCollision = true
	}

	// horizontal walls
	for y := 0; y <= m.Maze.Height(); y++ {
		j := float64(y * CellWidth)
		for x := 0; x < m.Maze.Width(); x++ {
			if *m.Maze.WallAbove(x, y) {
				i := float64(x * CellWidth)
				centerX, centerY = i+GlobalXOffset+CellWidth/2-tankCenterX, j+GlobalYOffset-tankCenterY
				if Judge() {
					HorizontalCollision = true
				}
			}
		}
	}

	// vertical walls
	for x := 0; x <= m.Maze.Width(); x++ {
		i := float64(x * CellWidth)
		for y := 0; y < m.Maze.Height(); y++ {
			if *m.Maze.WallLeftOf(x, y) {
				j := float64(y * CellWidth)
				centerX, centerY = i+GlobalXOffset-tankCenterX, j+GlobalYOffset-tankCenterY+CellWidth/2
				if Judge2() {
					VerticalCollision = true
				}
			}
		}
	}
	t.TankAHorizontalCollision, t.TankAVerticalCollision = HorizontalCollision, VerticalCollision
}

func BulletHitTank(t *Tank, bullet *Bullet) bool {
	if !bullet.Valid || bullet.JustShoot < 8 {
		return false
	}

	tankCenterX, tankCenterY := t.GetX(), t.GetY()

	axis1First, axis1Second := math.Cos(t.Angle), -math.Sin(t.Angle)
	axis2First, axis2Second := math.Sin(t.Angle), math.Cos(t.Angle)

	centerX, centerY := bullet.X-tankCenterX, bullet.Y-tankCenterY

	var Judge func() bool

	Judge = func() bool {
		projV := math.Abs(centerX*axis1First + centerY*axis1Second)
		proRadius := math.Abs(axis1First)*2 + math.Abs(axis1Second)*2
		if proRadius+TankHeight/2 <= projV {
			return false
		}
		projV = math.Abs(centerX*axis2First + centerY*axis2Second)
		proRadius = math.Abs(axis2First)*2 + math.Abs(axis2Second)*2
		if proRadius+TankWidth/2 <= projV {
			return false
		}
		projV = math.Abs(centerX)
		proRadius = math.Abs(axis1First)*TankHeight/2 + math.Abs(axis2First)*TankWidth/2
		if proRadius+2 <= projV {
			return false
		}
		projV = math.Abs(centerY)
		proRadius = math.Abs(axis1Second)*TankHeight/2 + math.Abs(axis2Second)*TankWidth/2
		if proRadius+2 <= projV {
			return false
		}
		return true
	}

	return Judge()
}
