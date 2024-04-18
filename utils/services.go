package utils

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

//Définition des fonctions de service et de formattage des données

func MessagePixelToString(pixel MessagePixel) string {
	return sepM + sepP + "positionX" + sepP + strconv.Itoa(pixel.PositionX) + sepM + sepP + "positionY" + sepP +
		strconv.Itoa(pixel.PositionY) + sepM + sepP + "R" + sepP + strconv.Itoa(pixel.Rouge) + sepM + sepP + "G" +
		sepP + strconv.Itoa(pixel.Vert) + sepM + sepP + "B" + sepP + strconv.Itoa(pixel.Bleu)
}

func MessageToString(message Message) string {
	c := ""
	if message.Couleur {
		c = "jaune"
	} else {
		c = "blanc"
	}
	return MessagePixelToString(message.Pixel) + sepM + sepP + "horloge" + sepP + strconv.Itoa(message.Horloge) +
		sepM + sepP + "nom" + sepP + message.Nom + sepM + sepP + "couleur" + sepP + c +
		sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.Prepost)

}

func MessageEtatToString(etat MessageEtat) string {
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
		l += MessagePixelToString(pixel)
	}

	return sep1 + sep2 + "etat" + sep2 + l + sep1 + sep2 + "bilan" + sep2 + strconv.Itoa(etat.Bilan)
}

func TrouverValeur(message string, cle string) string {
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
	posX, _ := strconv.Atoi(TrouverValeur(str, "positionX"))
	posY, _ := strconv.Atoi(TrouverValeur(str, "positionY"))
	r, _ := strconv.Atoi(TrouverValeur(str, "R"))
	v, _ := strconv.Atoi(TrouverValeur(str, "G"))
	b, _ := strconv.Atoi(TrouverValeur(str, "B"))

	messagepixel := MessagePixel{posX, posY, r, v, b}
	return messagepixel
}

func StringToMessage(str string) Message {
	messagepixel := StringToMessagePixel(str)
	h, _ := strconv.Atoi(TrouverValeur(str, "horloge"))
	n := TrouverValeur(str, "nom")
	cV := TrouverValeur(str, "couleur")
	var c Couleur
	if cV == "jaune" {
		c = Jaune
	} else {
		c = Blanc
	}
	prep, _ := strconv.ParseBool(TrouverValeur(str, "prepost"))
	message := Message{messagepixel, h, n, c, prep}
	return message
}

func StringToMessageEtat(str string) MessageEtat {
	var l list.List
	tabtousmesspix := TrouverValeur(str, "etat")
	fmt.Println(tabtousmesspix)
	tabtousmesspixsplit := strings.Split(tabtousmesspix, "_")
	for _, messpixel := range tabtousmesspixsplit {
		if messpixel != "" {
			fmt.Println(messpixel)
			l.PushBack(StringToMessagePixel(messpixel))
		}
	}
	b, _ := strconv.Atoi(TrouverValeur(str, "bilan"))
	messageetat := MessageEtat{l, b}
	return messageetat
}

func Recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}
