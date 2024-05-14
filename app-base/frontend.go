package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"utils"
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
	Matrix        [][]Pixel
	ColorWheel    *ebiten.Image
	SaveLogo      *ebiten.Image
	SelectedColor color.RGBA
}

func (g Game) UpdateMatrix(x int, y int, CR uint8, CG uint8, CB uint8) {
	g.Matrix[x][y] = Pixel{
		R: CR,
		G: CG,
		B: CB,
	}
}

func envoyerPixel(positionX int, positionY int, rouge int, vert int, bleu int) {
	messagePixel := utils.MessagePixel{positionX, positionY, rouge, vert, bleu}
	envoyerMessage(utils.MessagePixelToString(messagePixel))
}

func (g *Game) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= 80 && x <= 100 && y >= 80 && y <= 100 {
			if saveAccess {
				clicGaucheSaveLogo()
				saveAccess = false
			}
		} else if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
			clicGaucheMatrice(g, y, x, int(g.SelectedColor.R), int(g.SelectedColor.G), int(g.SelectedColor.B))
			// Les coordonnées de la souris ne sont pas comme est ordonnée la matrice, donc inversion de x et de y
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
			x_pourc := x * 100 / 50
			y_pourc := y * 100 / 50
			R, G, B, _ := g.ColorWheel.At(x_pourc-100, y_pourc).RGBA()
			g.SelectedColor = color.RGBA{uint8(R), uint8(G), uint8(B), 0xFF}
		}
	}
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	// Draw the main image
	screen.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	// Adjust position based on desired layout
	op.GeoM.Translate(0, 0)

	img := createImageFromMatrix(g.Matrix)
	screen.DrawImage(img, op)

	// Draw the color wheel
	colorWheelOp := &ebiten.DrawImageOptions{}
	colorWheelOp.GeoM.Translate(100, 0)
	colorWheelOp.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(g.ColorWheel, colorWheelOp)

	//Draw the save button
	saveLogoOp := &ebiten.DrawImageOptions{}
	saveLogoOp.GeoM.Translate(280, 280)
	saveLogoOp.GeoM.Scale(0.3, 0.3)
	screen.DrawImage(g.SaveLogo, saveLogoOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return len(g.Matrix[0]), len(g.Matrix)
}
