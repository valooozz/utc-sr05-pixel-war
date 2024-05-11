package utils

import (
	"strconv"
	"strings"
)

///////////////
// MessagePixel
///////////////

func MessagePixelToString(pixel MessagePixel) string {
	return sepM + sepP + "positionX" + sepP + strconv.Itoa(pixel.PositionX) + sepM + sepP + "positionY" + sepP +
		strconv.Itoa(pixel.PositionY) + sepM + sepP + "R" + sepP + strconv.Itoa(pixel.Rouge) + sepM + sepP + "G" +
		sepP + strconv.Itoa(pixel.Vert) + sepM + sepP + "B" + sepP + strconv.Itoa(pixel.Bleu)
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

//////////
// Message
//////////

func MessageToString(message Message) string {
	c := ""
	if message.Couleur {
		c = "jaune"
	} else {
		c = "blanc"
	}
	return MessagePixelToString(message.Pixel) + sepM + sepP + "horloge" + sepP + strconv.Itoa(message.Horloge) +
		sepM + sepP + "vectorielle" + sepP + HorlogeVectorielleToString(message.Vectorielle) + sepM + sepP + "nom" + sepP + message.Nom + sepM + sepP + "couleur" + sepP + c +
		sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.Prepost)

}

func StringToMessage(str string) Message {
	messagepixel := StringToMessagePixel(str)
	h, _ := strconv.Atoi(TrouverValeur(str, "horloge"))
	hv := TrouverValeur(str, "vectorielle")
	n := TrouverValeur(str, "nom")
	cV := TrouverValeur(str, "couleur")
	var c Couleur
	if cV == "jaune" {
		c = Jaune
	} else {
		c = Blanc
	}
	prep, _ := strconv.ParseBool(TrouverValeur(str, "prepost"))
	message := Message{messagepixel, h, StringToHorlogeVectorielle(hv), n, c, prep}
	return message
}

//////////////
// MessageEtat
//////////////

func MessageEtatToString(etat MessageEtat) string {
	sep1 := "~"
	sep2 := ","
	return sep1 + sep2 + "etat" + sep2 + EtatLocalToString(etat.EtatLocal)
}

func StringToMessageEtat(str string) MessageEtat {
	etatLocal := StringToEtatLocal(TrouverValeur(str, "etat"))

	return MessageEtat{etatLocal}
}

////////////
// EtatLocal
////////////

func EtatLocalToString(etatLocal EtatLocal) string {
	sep1 := "#"
	sep2 := ";"
	l := ""
	for _, messagePixel := range etatLocal.ListMessagePixel {
		l += "_"
		l += MessagePixelToString(messagePixel)
	}

	return sep1 + sep2 + "nom" + sep2 + etatLocal.NomSite +
		sep1 + sep2 + "vectorielle" + sep2 + HorlogeVectorielleToString(etatLocal.Vectorielle) +
		sep1 + sep2 + "liste" + sep2 + l
}

func StringToEtatLocal(str string) EtatLocal {
	var liste []MessagePixel
	listeMessagePixel := TrouverValeur(str, "liste")
	strVectorielle := TrouverValeur(str, "vectorielle")
	tabListeMessagePixel := strings.Split(listeMessagePixel, "_")

	for _, strMessagePixel := range tabListeMessagePixel {
		if strMessagePixel != "" {
			liste = append(liste, StringToMessagePixel(strMessagePixel))
		}
	}

	return EtatLocal{TrouverValeur(str, "nom"), StringToHorlogeVectorielle(strVectorielle), liste}
}

/////////////////////
// HorlogeVectorielle
/////////////////////

func HorlogeVectorielleToString(horloge HorlogeVectorielle) string {
	sep1 := "_"
	sep2 := ":"
	str := ""

	for site := range horloge {
		str += sep1
		str += site
		str += sep2
		str += strconv.Itoa(horloge[site])
	}

	return str
}

func StringToHorlogeVectorielle(str string) HorlogeVectorielle {
	horloge := HorlogeVectorielle{}
	listeSites := strings.Split(str, "_")

	for _, strSite := range listeSites {
		if strSite != "" {
			hSite := strings.Split(strSite, ":")
			horloge[hSite[0]], _ = strconv.Atoi(hSite[1])
		}
	}

	return horloge
}
