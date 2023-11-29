package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	Red   = color.RGBA{R: 0xff, G: 0x35, B: 0x00, A: 0xff}
	Green = color.RGBA{R: 0x20, G: 0xf0, B: 0x11, A: 0xff}
)

// Mode 定义游戏状态
type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeAWin
	ModeBWin
)

const TankLineSpeed = 2.4
const TankAngularSpeed = math.Pi / 50

type Game struct {
	// Mode 当前的模式
	Mode Mode
	// 玩家 A 胜利局数
	aCount int
	// 玩家 B 胜利局数
	bCount int

	TankA *Tank
	TankB *Tank

	GameMap    *Map
	pressedKey []ebiten.Key

	// 对局结束后的暂停
	breaking int
}

const (
	titleFontSize = fontSize * 2
	fontSize      = 32
	smallFontSize = fontSize / 2
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	switch g.Mode {
	case ModeTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
			g.Mode = ModeGame
		}
	case ModeGame:
		g.pressedKey = inpututil.AppendPressedKeys(g.pressedKey[:0])

		for _, key := range g.pressedKey {
			if key == ebiten.KeyArrowUp {
				g.TankA.Vx -= TankLineSpeed * math.Cos(g.TankA.Angle)
				g.TankA.Vy -= TankLineSpeed * math.Sin(g.TankA.Angle)
			} else if key == ebiten.KeyArrowDown {
				g.TankA.Vx += TankLineSpeed * math.Cos(g.TankA.Angle)
				g.TankA.Vy += TankLineSpeed * math.Sin(g.TankA.Angle)
			} else if key == ebiten.KeyArrowLeft {
				g.TankA.Angle -= TankAngularSpeed
			} else if key == ebiten.KeyArrowRight {
				g.TankA.Angle += TankAngularSpeed
			} else if key == ebiten.KeyW {
				g.TankB.Vx -= TankLineSpeed * math.Cos(g.TankB.Angle)
				g.TankB.Vy -= TankLineSpeed * math.Sin(g.TankB.Angle)
			} else if key == ebiten.KeyS {
				g.TankB.Vx += TankLineSpeed * math.Cos(g.TankB.Angle)
				g.TankB.Vy += TankLineSpeed * math.Sin(g.TankB.Angle)
			} else if key == ebiten.KeyA {
				g.TankB.Angle -= TankAngularSpeed
			} else if key == ebiten.KeyD {
				g.TankB.Angle += TankAngularSpeed
			}
		}
		g.GameMap.CheckCollision(g.TankA)
		g.GameMap.CheckCollision(g.TankB)
		g.TankA.UpdateSpeed()
		g.TankB.UpdateSpeed()
		if repeatingKeyPressed(ebiten.KeyJ) {
			for _, bullet := range g.TankA.Bullets {
				if !bullet.Valid {
					g.TankA.Fire(bullet)
					break
				}
			}
		}
		for _, bullet := range g.TankA.Bullets {
			bullet.UpdateBulletPos()
			if BulletHitTank(g.TankA, bullet) {
				g.Mode = ModeBWin
			}
			if BulletHitTank(g.TankB, bullet) {
				g.Mode = ModeAWin
			}
			g.GameMap.Hit(bullet, 0, 0)
		}
		if repeatingKeyPressed(ebiten.KeySpace) {
			for _, bullet := range g.TankB.Bullets {
				if !bullet.Valid {
					g.TankB.Fire(bullet)
					break
				}
			}
		}
		for _, bullet := range g.TankB.Bullets {
			bullet.UpdateBulletPos()
			if BulletHitTank(g.TankB, bullet) {
				g.Mode = ModeAWin
			}
			if BulletHitTank(g.TankA, bullet) {
				g.Mode = ModeBWin
			}
			g.GameMap.Hit(bullet, 0, 0)
		}
	case ModeAWin:
		g.breaking++
		if g.breaking > 120 {
			g.aCount++
			g.breaking = 0
			g.Mode = ModeGame
			g.GameMap = NewMap()
			g.TankA = NewTank(colornames.Grey, 150, 100, math.Pi)
			g.TankB = NewTank(Red, 810, 520, 0)
		}
	case ModeBWin:
		g.breaking++
		if g.breaking > 120 {
			g.bCount++
			g.breaking = 0
			g.Mode = ModeGame
			g.GameMap = NewMap()
			g.TankA = NewTank(colornames.Grey, 150, 100, math.Pi)
			g.TankB = NewTank(Red, 810, 520, 0)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)
	ebitenutil.DebugPrint(screen, "FPS: "+strconv.FormatFloat(ebiten.ActualFPS(), 'f', 1, 64))
	ebitenutil.DebugPrint(screen, "\nTPS: "+strconv.FormatFloat(ebiten.ActualTPS(), 'f', 1, 64))

	switch g.Mode {
	case ModeTitle:
		titleTexts := []string{"Tank Trouble"}
		texts := []string{"Press F or Space To Start."}
		name := string("by bz2021")
		for i, l := range titleTexts {
			x := (960 - len(l)*titleFontSize) / 2
			text.Draw(screen, l, titleArcadeFont, x-2, (i+3)*titleFontSize-2, color.Black)
		}
		for i, l := range texts {
			x := (960 - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x-2, (i+4)*fontSize+298, color.Black)
		}
		for i, l := range titleTexts {
			x := (960 - len(l)*titleFontSize) / 2
			text.Draw(screen, l, titleArcadeFont, x, (i+3)*titleFontSize, color.White)
		}
		for i, l := range texts {
			x := (960 - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize+300, color.White)
		}
		text.Draw(screen, name, smallArcadeFont, 390, 620, color.Black)
		text.Draw(screen, name, smallArcadeFont, 391, 621, color.White)
	case ModeGame:
		g.GameMap.DrawMap(screen)
		g.TankA.DrawTank(screen)
		g.TankB.DrawTank(screen)
		for _, b := range g.TankA.Bullets {
			b.DrawBullet(screen)
		}
		for _, b := range g.TankB.Bullets {
			b.DrawBullet(screen)
		}
		aWin := "Grey: " + strconv.Itoa(g.aCount)
		bWin := "Red: " + strconv.Itoa(g.bCount)
		text.Draw(screen, aWin, smallArcadeFont, 120, 40, color.Black)
		text.Draw(screen, bWin, smallArcadeFont, 550, 40, color.Black)
		//g.TankA.laser.DrawLaser(screen, g.TankA)
	case ModeAWin:
		g.GameMap.DrawMap(screen)
		g.TankA.DrawTank(screen)
		if g.breaking/10%3 == 0 {
			newTank := NewTank(colornames.White, g.TankB.GetX(), g.TankB.GetY(), g.TankB.Angle)
			newTank.DrawTank(screen)
		} else {
			g.TankB.DrawTank(screen)
		}
		for _, b := range g.TankA.Bullets {
			b.DrawBullet(screen)
		}
		for _, b := range g.TankB.Bullets {
			b.DrawBullet(screen)
		}
		aWin := "Grey: " + strconv.Itoa(g.aCount)
		bWin := "Red: " + strconv.Itoa(g.bCount)
		text.Draw(screen, aWin, smallArcadeFont, 120, 40, color.Black)
		text.Draw(screen, bWin, smallArcadeFont, 550, 40, color.Black)
	case ModeBWin:
		g.GameMap.DrawMap(screen)
		g.TankB.DrawTank(screen)
		if g.breaking/10%3 == 0 {
			newTank := NewTank(colornames.White, g.TankA.GetX(), g.TankA.GetY(), g.TankA.Angle)
			newTank.DrawTank(screen)
		} else {
			g.TankA.DrawTank(screen)
		}
		for _, b := range g.TankA.Bullets {
			b.DrawBullet(screen)
		}
		for _, b := range g.TankB.Bullets {
			b.DrawBullet(screen)
		}
		aWin := "Grey: " + strconv.Itoa(g.aCount)
		bWin := "Red: " + strconv.Itoa(g.bCount)
		text.Draw(screen, aWin, smallArcadeFont, 120, 40, color.Black)
		text.Draw(screen, bWin, smallArcadeFont, 550, 40, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 960, 640
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 20
		interval = 6
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
