package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math"
	"os"
	"time"
)

var (
	nScreenWidth, nScreenHeight int = 120, 40

	fPlayerX float64 = 1.0
	fPlayerY float64 = 1.0
	fPlayerA float64 = 0.0

	nMapHeight int = 16
	nMapWidth  int = 16

	fFov   float64 = 3.14159 / 4
	fDepth float64 = 16.0
)

func printMap(m []rune) {
	for x := 0; x < nMapWidth; x++ {
		for y := 0; y < nMapHeight; y++ {
			if x == int(fPlayerX) && y == int(fPlayerY) {
				termbox.SetCell(x, y, 'P', termbox.ColorWhite, termbox.ColorBlack)
			} else {
				termbox.SetCell(x, y, m[y*nMapHeight+x], termbox.ColorWhite, termbox.ColorBlack)
			}
		}
	}
}

func setCell(x, y int, ch rune, color termbox.Attribute) {
	termbox.SetCell(x, y, ch, color, termbox.ColorBlack)
}

func main() {
	err := termbox.Init()
	defer termbox.Close()
	if err != nil {
		fmt.Println("Init error")
		os.Exit(1)
	}

	gameMap := []rune(
		"################" +
			"#..............#" +
			"#..............#" +
			"#....###########" +
			"#..............#" +
			"#..............#" +
			"###########....#" +
			"#..............#" +
			"#..............#" +
			"#.....##########" +
			"#..............#" +
			"#..............#" +
			"############...#" +
			"#..............#" +
			"#..............#" +
			"################")

	tp1 := time.Now()
	tp2 := time.Now()

	// Game loop
	for {
		tp2 = time.Now()
		elapsedTime := tp2.Sub(tp1)
		tp1 = tp2
		fElapsedTime := math.Min(elapsedTime.Seconds(), 0.5)

		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyArrowLeft:
				fPlayerA -= 0.1
			case termbox.KeyArrowRight:
				fPlayerA += 0.1
			case termbox.KeyArrowUp:
				fPlayerX += math.Sin(fPlayerA) * 2.0 * fElapsedTime
				fPlayerY += math.Cos(fPlayerA) * 2.0 * fElapsedTime

				if gameMap[int(fPlayerY)*nMapHeight+int(fPlayerX)] == '#' {
					fPlayerX -= math.Sin(fPlayerA) * 2.0 * fElapsedTime
					fPlayerY -= math.Cos(fPlayerA) * 2.0 * fElapsedTime
				}
			case termbox.KeyArrowDown:
				fPlayerX -= math.Sin(fPlayerA) * 2.0 * fElapsedTime
				fPlayerY -= math.Cos(fPlayerA) * 2.0 * fElapsedTime

				if gameMap[int(fPlayerY)*nMapHeight+int(fPlayerX)] == '#' {
					fPlayerX += math.Sin(fPlayerA) * 2.0 * fElapsedTime
					fPlayerY += math.Cos(fPlayerA) * 2.0 * fElapsedTime
				}
			}

		case termbox.EventError:
			os.Exit(1)
		}

		for x := 0; x < nScreenWidth; x++ {
			var fRayAngle float64 = (fPlayerA - fFov/2.0) + (float64(x) / float64(nScreenWidth) * fFov)

			var fDistanceToWall float64 = 0
			bHitWall := false

			fEyeX := math.Sin(fRayAngle)
			fEyeY := math.Cos(fRayAngle)

			for !bHitWall && fDistanceToWall < fDepth {

				fDistanceToWall += 0.1

				nTestX := int(fPlayerX + fEyeX*fDistanceToWall)
				nTestY := int(fPlayerY + fEyeY*fDistanceToWall)

				if nTestX < 0 || nTestX >= nMapWidth || nTestY < 0 || nTestY >= nMapHeight {
					bHitWall = true
					fDistanceToWall = fDepth
				} else {
					if gameMap[nTestY*nMapWidth+nTestX] == '#' {
						bHitWall = true
					}
				}
			}

			nCeiling := int(float64(nScreenHeight)/2.0 - float64(nScreenHeight)/fDistanceToWall)
			nFloor := nScreenHeight - nCeiling

			var nShade rune

			if fDistanceToWall <= fDepth/4.0 {
				nShade = '█'
			} else if fDistanceToWall < fDepth/3.0 {
				nShade = '┼'
			} else if fDistanceToWall < fDepth/2.0 {
				nShade = '┬'
			} else if fDistanceToWall < fDepth {
				nShade = '─'
			} else {
				nShade = ' '
			}

			for y := 0; y < nScreenHeight; y++ {
				if y < nCeiling {
					setCell(x, y, ' ', termbox.ColorBlack)
				} else if y > nCeiling && y <= nFloor {
					setCell(x, y, nShade, termbox.ColorWhite)
				} else {
					b := (float64(y) - float64(nScreenHeight)/2.0) / (float64(nScreenHeight) / 2.0)
					var floorShade rune
					if b < 0.25 {
						floorShade = ' '
					} else if b < 0.5 {
						floorShade = 'x'
					} else if b < 0.75 {
						floorShade = '`'
					} else if b < 0.9 {
						floorShade = '.'
					} else {
						floorShade = ' '
					}
					setCell(x, y, floorShade, termbox.ColorBlue)
				}
			}
		}
		printMap(gameMap)
		err := termbox.Flush()
		if err != nil {
			os.Exit(1)
		}
	}
}
