package utils

import (
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
		sepM + sepP + "vectorielle" + sepP + HorlogeVectorielleToString(message.Vectorielle) + sepM + sepP + "nom" + sepP + message.Nom + sepM + sepP + "couleur" + sepP + c +
		sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.Prepost)

}

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

func MessageEtatToString(etat MessageEtat) string {
	sep1 := "~"
	sep2 := ","
	return sep1 + sep2 + "etat" + sep2 + EtatLocalToString(etat.EtatLocal)
}

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

func StringToMessageEtat(str string) MessageEtat {
	etatLocal := StringToEtatLocal(TrouverValeur(str, "etat"))

	return MessageEtat{etatLocal}
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

func MajHorlogeVectorielle(monNom string, locale, recue HorlogeVectorielle) HorlogeVectorielle {

	// On met à jour les champs présents dans l'horloge locale
	for site, valeurLocale := range locale {
		valeurRecue, ok := recue[site]
		if ok {
			if valeurRecue > valeurLocale {
				locale[site] = valeurRecue
			}
			delete(recue, site)
		}
	}

	// On ajoute les champs restants dans l'horloge reçue
	for site, valeurRecue := range recue {
		locale[site] = valeurRecue
	}

	// On incrémente l'horloge du site local
	locale[monNom]++

	return locale
}

func CopyHorlogeVectorielle(horlogeVectorielle HorlogeVectorielle) HorlogeVectorielle {

	var copie = HorlogeVectorielle{}

	for key, val := range horlogeVectorielle {
		copie[key] = val
	}

	return copie
}

func CoupureEstCoherente(etatGlobal EtatGlobal) bool {
	isProcessed := make(map[string]bool)
	mapMax := make(map[string]int)

	// Initialisation
	for _, etatLocal := range etatGlobal.ListEtatLocal {
		for site, _ := range etatLocal.Vectorielle {
			isProcessed[site] = false
			mapMax[site] = 0
		}
	}

	for _, etatLocal := range etatGlobal.ListEtatLocal {
		for site, horloge := range etatLocal.Vectorielle {
			if mapMax[site] < horloge { // Si l'horloge est plus grande que le max enregistré
				if isProcessed[site] { // Si on a déjà passé le site, la coupure n'est pas cohérente
					return false
				} else { // Sinon, on met à jour le max
					mapMax[site] = horloge
				}
			} else if mapMax[site] > horloge && etatLocal.NomSite == site {
				return false // Si le max du site est plus grand que l'horloge de ce site sur ce site, la coupure n'est pas cohérente
			}
		}
		isProcessed[etatLocal.NomSite] = true // Le site a été process
	}

	return true
}

func MajEtatLocal(etatLocal EtatLocal, newMessagePixel MessagePixel) EtatLocal {
	var found = false
	for i, pixel := range etatLocal.ListMessagePixel {
		if pixel.PositionX == newMessagePixel.PositionX && pixel.PositionY == newMessagePixel.PositionY {
			pixel.Rouge = newMessagePixel.Rouge
			pixel.Vert = newMessagePixel.Vert
			pixel.Bleu = newMessagePixel.Bleu
			etatLocal.ListMessagePixel[i] = pixel
			found = true
		}
	}

	if !found {
		etatLocal.ListMessagePixel = append(etatLocal.ListMessagePixel, newMessagePixel)
	}

	return etatLocal
}

func CopyEtatLocal(etatLocal EtatLocal) EtatLocal {
	var copie = EtatLocal{
		NomSite:          etatLocal.NomSite,
		Vectorielle:      etatLocal.Vectorielle,
		ListMessagePixel: []MessagePixel{},
	}

	for _, mp := range etatLocal.ListMessagePixel {
		copie.ListMessagePixel = append(copie.ListMessagePixel, mp)
	}

	return copie
}

func Recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}
