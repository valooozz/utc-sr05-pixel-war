package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"utils"
)

// Structure Pixel composé de trois entiers non signés sur 8 bits pour la valeur RGB
type Pixel struct {
	R, G, B uint8
}

// Fonction pour initialiser l'image de la matrice
func createImageFromMatrix(matrix [][]Pixel) *ebiten.Image {
	// Initialisation d'une image ed taille 50x50
	width := 50
	height := 50
	img := ebiten.NewImage(width, height)
	// On initialise chaque pixel de l'image en blanc (la matrice ne contient que des pixels blancs initialement
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := matrix[y][x]
			img.Set(x, y, color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: 0xFF})
		}
	}
	return img
}

// Définition du type Game, composé d'une matrice de Pixel, d'une roue de couleur, d'un logo de sauvegarde
// ainsi qu'une couleur selectionnée
type Game struct {
	Matrix        [][]Pixel
	ColorWheel    *ebiten.Image
	SaveLogo      *ebiten.Image
	SelectedColor color.RGBA
}

func (g Game) UpdateMatrix(x int, y int, CR uint8, CG uint8, CB uint8) {
	// On met à jour le pixel de la matrice
	g.Matrix[x][y] = Pixel{
		R: CR,
		G: CG,
		B: CB,
	}
}

func envoyerPixel(positionX int, positionY int, rouge int, vert int, bleu int) {
	// On envoie un message contenant le pixel posé à l'app de contrôle
	messagePixel := utils.MessagePixel{positionX, positionY, rouge, vert, bleu}
	envoyerMessage(utils.MessagePixelToString(messagePixel))
}

// Update met à jour l'état du jeu
// Update est appelée à chaque tick (1/60 [s] par défaut)
func (g *Game) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()
	// Action de clic gauche
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		// Si le curseur se situe sur le logo de sauvegarde
		if x >= 80 && x <= 100 && y >= 80 && y <= 100 {
			if saveAccess {
				clicGaucheSaveLogo()
				saveAccess = false
			}
			// Si le curseur se situe dans la zone de dessin
		} else if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
			clicGaucheMatrice(g, y, x, int(g.SelectedColor.R), int(g.SelectedColor.G), int(g.SelectedColor.B))
			// Les coordonnées récupérées par le curseur ne sont pas ordonnées de la même manière qu'est ordonnée
			// la matrice, il faut inverser x et y
		}
	}
	// Action de clic droit
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		// Lors de l'exécution la roue de couleur est traitée comme occupant tout l'espace de l'interface
		// Il est donc nécessaire de redimensionner les coordonnées pour s'adpater à sa réelle position
		if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
			x_pourc := x * 100 / 50
			y_pourc := y * 100 / 50
			// On sauvegarde la couleur sélectionnée
			R, G, B, _ := g.ColorWheel.At(x_pourc-100, y_pourc).RGBA()
			g.SelectedColor = color.RGBA{uint8(R), uint8(G), uint8(B), 0xFF}
		}
	}
	return nil
}

// Draw met à jour le visuel du jeu
// Elle est appelée à chaque frame (1/60 [s] sur un écran 60Hz)
func (g Game) Draw(screen *ebiten.Image) {

	// On dessine l'image de la matrice
	screen.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	img := createImageFromMatrix(g.Matrix)
	// On dessine l'image de la matrice avec les options op (image situé en haut à gauche)
	screen.DrawImage(img, op)

	// On dessine la roue de couleur
	colorWheelOp := &ebiten.DrawImageOptions{}
	colorWheelOp.GeoM.Translate(100, 0)
	colorWheelOp.GeoM.Scale(0.5, 0.5)
	// On dessine la roue de couleur avec les options colorWheelOp (translationde 100 sur l'axe x
	// et réduction de la taille de l'image par moitié)
	screen.DrawImage(g.ColorWheel, colorWheelOp)

	// On dessine le bouton de sauvegarde
	saveLogoOp := &ebiten.DrawImageOptions{}
	saveLogoOp.GeoM.Translate(280, 280)
	saveLogoOp.GeoM.Scale(0.3, 0.3)
	// On dessine le bouton de sauvegarde avec les options saveLogoOp (translationde 280 sur l'axe x et y
	// et réduction de la taille de l'image à un tier de sa taille)
	screen.DrawImage(g.SaveLogo, saveLogoOp)
}

// La documentation pour Layout explique
// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return len(g.Matrix[0]), len(g.Matrix)
}
