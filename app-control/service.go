package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

//Définition des fonctions de service et de formattage des données

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
		sepM + sepP + "nom" + sepP + message.nom + sepM + sepP + "couleur" + sepP + c +
		sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.prepost)

}

func messageEtatToString(etat MessageEtat) string {
	sep1 := "~"
	sep2 := ","
	l := ""
	for e := etat.EG.Front(); e != nil; e = e.Next() {
		l += "_"
		pixel, ok := e.Value.(MessagePixel)
		if !ok {
			fmt.Println("Conversion to MessagePixel failed")
			return ""
		}
		l += messagePixelToString(pixel)
	}

	return sep1 + sep2 + "etat" + sep2 + l + sep1 + sep2 + "bilan" + sep2 + strconv.Itoa(etat.bilan)
}

func trouverValeur(message string, cle string) string {
	if len(message) < 4 {
		return ""
	}
	sep := message[0:1]
	tabToutesCleValeur := strings.Split(message[1:], sep)
	for _, cleV := range tabToutesCleValeur {
		equ := cleV[0:1]
		tabCleValeur := strings.Split(cleV[1:], equ)
		if tabCleValeur[0] == cle {
			return tabCleValeur[1]
		}
	}
	return ""
}

func StringToMessagePixel(str string) MessagePixel {
	posX, _ := strconv.Atoi(trouverValeur(str, "positionX"))
	posY, _ := strconv.Atoi(trouverValeur(str, "positionY"))
	r, _ := strconv.Atoi(trouverValeur(str, "R"))
	v, _ := strconv.Atoi(trouverValeur(str, "G"))
	b, _ := strconv.Atoi(trouverValeur(str, "B"))

	messagepixel := MessagePixel{posX, posY, r, v, b}
	return messagepixel
}

func StringToMessage(str string) Message {
	messagepixel := StringToMessagePixel(str)
	h, _ := strconv.Atoi(trouverValeur(str, "horloge"))
	n := trouverValeur(str, "nom")
	cV := trouverValeur(str, "couleur")
	var c Couleur
	if cV == "jaune" {
		c = Jaune
	} else {
		c = Blanc
	}
	prep, _ := strconv.ParseBool(trouverValeur(str, "prepost"))
	message := Message{messagepixel, h, n, c, prep}
	return message
}

func StringToMessageEtat(str string) MessageEtat {
	var l list.List
	tabtousmesspix := trouverValeur(str, "etat")
	fmt.Println(tabtousmesspix)
	tabtousmesspixsplit := strings.Split(tabtousmesspix, "_")
	for _, messpixel := range tabtousmesspixsplit {
		if messpixel != "" {
			fmt.Println(messpixel)
			l.PushBack(StringToMessagePixel(messpixel))
		}
	}
	b, _ := strconv.Atoi(trouverValeur(str, "bilan"))
	messageetat := MessageEtat{l, b}
	return messageetat
}

func recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}
