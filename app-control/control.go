package main

import (
	"container/list"
	"flag"
	"os"
	"strconv"
	"sync"
)

// Définition des variables
var mutex = &sync.Mutex{}
var K = 0
var maCouleur = Blanc
var initiateur = false

var pNom = flag.String("n", "controle", "nom")
var nom string

// Définition des constantes
// tailleMap = ?
// nbSites = ?
const sepM = "/" //séparateur dans les messages
const sepP = "=" //séparateur ddans les paires clé/valeur
var rouge string = "\033[1;31m"
var orange string = "\033[1;33m"
var raz string = "\033[0;00m"

// Définition des types
type Couleur bool

const (
	Blanc Couleur = false
	Jaune Couleur = true
)

type MessagePixel struct {
	positionX int
	positionY int
	rouge     int
	vert      int
	bleu      int
}

type Message struct {
	pixel   MessagePixel
	horloge int
	nom     string
	couleur Couleur
	prepost bool //false pour les messages normaux
}

type EtatGlobal list.List //Sous-entendu une liste de MessagePixel

type MessageEtat struct {
	EG    list.List
	bilan int
}

func main() {
	flag.Parse()
	nom = *pNom + "-" + strconv.Itoa(os.Getpid())

	//go lecture()
	//go ecriture()
	//for {
	//	time.Sleep(time.Duration(60) * time.Second)
	//} // Pour attendre la fin des goroutines...
}
