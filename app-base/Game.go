package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Pixel struct {
	R, G, B uint8
}

func createImageFromMatrix(matrix [][]Pixel) *ebiten.Image {
	width := len(matrix[0])
	height := len(matrix)
	img := ebiten.NewImage(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := matrix[y][x]
			img.Set(x, y, color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: 0xFF})
		}
	}

	return img
}

type Game struct {
	matrix [][]Pixel
}

func (g *Game) UpdateMatrix(x int, y int, CR uint8, CG uint8, CB uint8) {
	g.matrix[x][y] = Pixel{
		R: CR,
		G: CG,
		B: CB,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	img := createImageFromMatrix(g.matrix)
	screen.DrawImage(img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return len(g.matrix[0]), len(g.matrix)
}

func main() {
	matrix := make([][]Pixel, 100)
	for y := 0; y < 100; y++ {
		matrix[y] = make([]Pixel, 100)
		for x := 0; x < 100; x++ {
			matrix[y][x] = Pixel{
				R: 255,
				G: 255,
				B: 255,
			}
		}
	}

	game := &Game{
		matrix: matrix,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Pixel-War")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, " ")
			if len(parts) == 5 {
				x, err := strconv.Atoi(parts[0])
				if err != nil {
					continue
				}
				y, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}
				cr, err := strconv.Atoi(parts[2])
				if err != nil {
					continue
				}
				cg, err := strconv.Atoi(parts[3])
				if err != nil {
					continue
				}
				cb, err := strconv.Atoi(parts[4])
				if err != nil {
					continue
				}
				game.UpdateMatrix(x, y, uint8(cr), uint8(cg), uint8(cb))
				fmt.Printf("Updated pixel at (%d, %d) to (%d, %d, %d)\n", x, y, cr, cg, cb)
			}
		}
	}()

	ebiten.RunGame(game)

}
