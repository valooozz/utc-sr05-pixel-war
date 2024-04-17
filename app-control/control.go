package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Définition des constantes
// tailleMap = ?
// nbSites = ?
const sepM = "/" //séparateur dans les messages
const sepP = "=" //séparateur ddans les paires clé/valeur

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
	//Ajouter le pid ?
	couleur Couleur
	prepost bool //false pour les messages normaux
}

type EtatGlobal list.List //Sous-entendu une liste de MessagePixel

type MessageEtat struct {
	EG    list.List
	bilan int
}

//Définition des fonctions de service

func messagePixelToString(pixel MessagePixel) string {
	return sepM + sepP + "positionX" + sepP + strconv.Itoa(pixel.positionX) + sepM + sepP + "positionY" + sepP +
		strconv.Itoa(pixel.positionY) + sepM + sepP + "R" + sepP + strconv.Itoa(pixel.rouge) + sepM + sepP + "G" +
		sepP + strconv.Itoa(pixel.vert) + sepM + sepP + "B" + sepP + strconv.Itoa(pixel.bleu)
}

func messageToString(message Message) string {
	c := ""
	if message.couleur {
		c = "jaune"
	} else {
		c = "blanc"
	}
	return messagePixelToString(message.pixel) + sepM + sepP + "horloge" + sepP + strconv.Itoa(message.horloge) +
		sepM + sepP + "couleur" + sepP + c + sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.prepost)

}

func messageEtatToString(etat MessageEtat) string {
	sep1 := "~"
	sep2 := ","
	l := ""
	for e := etat.EG.Front(); e != nil; e = e.Next() {
		l += messagePixelToString(e.Value)
	}

	return sep1 + sep2 + "etat" + sep2 + l + sep1 + sep2 + "bilan" + sep2 + strconv.Itoa(etat.bilan)
}

func trouverValeur(message string, cle string) string {
	if len(message) < 4 {
		return ""
	}
	sep := message[0:1]
	tab__toutes_cle_valeur := strings.Split(message[1:], sep)
	for _, cle_v := range tab__toutes_cle_valeur {
		equ := cle_v[0:1]
		tab_cle_valeur := strings.Split(cle_v[1:], equ)
		if tab_cle_valeur[0] == cle {
			return tab_cle_valeur[1]
		}
	}
	return ""
}

func StringToMessagePixel(str string) MessagePixel {
	pos_X, _ := strconv.Atoi(trouverValeur(str, "positionX"))
	pos_Y, _ := strconv.Atoi(trouverValeur(str, "positionY"))
	r, _ := strconv.Atoi(trouverValeur(str, "R"))
	v, _ := strconv.Atoi(trouverValeur(str, "G"))
	b, _ := strconv.Atoi(trouverValeur(str, "B"))

	messagepixel := MessagePixel{pos_X, pos_Y, r, v, b}
	return messagepixel
}

func StringToMessage(str string) Message {
	messagepixel := StringToMessagePixel(str)
	h, _ := strconv.Atoi(trouverValeur(str, "horloge"))
	c_v := trouverValeur(str, "couleur")
	var c Couleur
	if c_v == "jaune" {
		c = Jaune
	} else {
		c = Blanc
	}
	prep, _ := strconv.ParseBool(trouverValeur(str, "prepost"))
	message := Message{messagepixel, h, c, prep}
	return message
}

//func StringToMessageEtat(str string) MessageEtat {
//}

var mutex = &sync.Mutex{}
var K = 0
var maCouleur = Blanc

func main() {
	messagePixel := MessagePixel{positionX: 1, positionY: 2, rouge: 43, vert: 67, bleu: 98}
	fmt.Println(messagePixelToString(StringToMessagePixel(messagePixelToString(messagePixel))))
	fmt.Println("\n")
	message := Message{messagePixel, 0, Blanc, false}
	fmt.Println(messageToString(message))
	premier := trouverValeur(messagePixelToString(messagePixel), "positionX")
	fmt.Println(premier + "\n")
	deuxieme := trouverValeur(messagePixelToString(messagePixel), "positionZ")
	fmt.Println(deuxieme + "\n")

	//go lecture()
	//go ecriture()
	//for {
	//	time.Sleep(time.Duration(60) * time.Second)
	//} // Pour attendre la fin des goroutines...
}
