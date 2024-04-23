package main

import (
	"bufio"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	width := 50
	height := 50
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
	matrix        [][]Pixel
	colorWheel    *ebiten.Image
	selectedColor color.RGBA
}

func (g *Game) UpdateMatrix(x int, y int, CR uint8, CG uint8, CB uint8) {
	g.matrix[x][y] = Pixel{
		R: CR,
		G: CG,
		B: CB,
	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.UpdateMatrix(y, x, g.selectedColor.R, g.selectedColor.G, g.selectedColor.B)
		// Oui je sais c'est bizarre mais les coordonnées de la souris ne sont pas comme est ordonnée la matrice
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		x_pourc := x * 100 / 50
		y_pourc := y * 100 / 50

		R, G, B, _ := g.colorWheel.At(x_pourc-100, y_pourc).RGBA()
		g.selectedColor = color.RGBA{uint8(R), uint8(G), uint8(B), 0xFF}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the main image
	screen.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	// Adjust position based on desired layout (explained later)
	op.GeoM.Translate(0, 0)

	img := createImageFromMatrix(g.matrix)
	screen.DrawImage(img, op)

	// Draw the color wheel
	colorWheelOp := &ebiten.DrawImageOptions{}
	colorWheelOp.GeoM.Translate(100, 0)
	colorWheelOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(g.colorWheel, colorWheelOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return len(g.matrix[0]), len(g.matrix)
}

func main1() {
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
	colorWheel, _, err := ebitenutil.NewImageFromFile("color_wheel.png")
	if err != nil {
		panic(err)
	}

	game := &Game{
		matrix:        matrix,
		colorWheel:    colorWheel,
		selectedColor: color.RGBA{R: 0, G: 0, B: 0, A: 0xFF},
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
