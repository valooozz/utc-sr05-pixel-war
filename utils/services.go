package utils

import (
	"strings"
)

//Définition des fonctions de service et de formattage des données

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

func Recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
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
