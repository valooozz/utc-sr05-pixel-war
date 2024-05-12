package utils

import (
	"strconv"
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

func CoupureEstCoherente(etatGlobal EtatGlobal) (bool, map[string]int) {
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
					return false, mapMax
				} else { // Sinon, on met à jour le max
					mapMax[site] = horloge
				}
			} else if mapMax[site] > horloge && etatLocal.NomSite == site {
				return false, mapMax // Si le max du site est plus grand que l'horloge de ce site sur ce site, la coupure n'est pas cohérente
			}
		}
		isProcessed[etatLocal.NomSite] = true // Le site a été process
	}

	return true, mapMax
}

func MajEtatLocal(etatLocal EtatLocal, newMessagePixel MessagePixel) EtatLocal {
	var found = false
	for i, pixel := range etatLocal.ListMessagePixel {
		if memePosition(pixel, newMessagePixel) {
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

func ReconstituerCarte(etatGlobal EtatGlobal) []MessagePixel {
	var carte = etatGlobal.ListEtatLocal[0].ListMessagePixel
	var pixel MessagePixel

	for _, message := range etatGlobal.ListMessagePrepost {
		pixel = message.Pixel
		carte = replaceOrAddPixel(carte, pixel)
	}

	return carte
}

func ReconstituerCarteOld(etatGlobal EtatGlobal) []MessagePixel {
	var carte = etatGlobal.ListEtatLocal[0].ListMessagePixel

	for _, etatLocal := range etatGlobal.ListEtatLocal[1:] {
		for i, pixel := range etatLocal.ListMessagePixel {
			if i-1 > len(carte) {
				carte = append(carte, pixel)
				continue
			}
			if memePosition(carte[i], pixel) && !memeCouleur(carte[i], pixel) {
				found, prepost := getPrepostOnSamePosition(pixel, etatGlobal.ListMessagePrepost)
				if found {
					carte[i] = prepost
				}
			}
		}
	}

	return carte
}

/////////////////////
// Exclusion mutuelle
/////////////////////

// Tester et valider
func QuestionEntreeSC(site int, tabSC []MessageExclusionMutuelle) bool {
	cpt := 0

	if tabSC[site].Type != Requete {
		return false
	}

	for i := 0; i < len(tabSC); i++ {
		DisplayInfo("site"+strconv.Itoa(site), "Question", "tabSC["+strconv.Itoa(i)+"]="+strconv.Itoa(int(tabSC[i].Type))+" | "+strconv.Itoa(tabSC[i].Estampille.Horloge))
	}
	for numOtherSite := 0; numOtherSite < len(tabSC); numOtherSite++ {
		if numOtherSite == site {
			continue
		}

		if tabSC[site].Estampille.Horloge > tabSC[numOtherSite].Estampille.Horloge {
			return false
		}
		if tabSC[site].Estampille.Horloge < tabSC[numOtherSite].Estampille.Horloge {
			cpt++
		} else {
			if tabSC[site].Estampille.Site > tabSC[numOtherSite].Estampille.Site {
				return false
			}
			cpt++
		}
	}

	if cpt == len(tabSC)-1 {
		DisplayInfo("Coucou", "Question", "Accepté pour "+strconv.Itoa(site))
		return true
	}
	return false
}

// Tester et valider
func InitialisationNumSite(site string) int {
	StartNumberIndex := 1
	SiteString := site[StartNumberIndex:len(site)]
	NumSite, _ := strconv.Atoi(SiteString)
	return NumSite
}

////////////////////
// FONCTIONS PRIVEES
////////////////////

func memePosition(pixel1, pixel2 MessagePixel) bool {
	if pixel1.PositionX == pixel2.PositionX && pixel1.PositionY == pixel2.PositionY {
		return true
	}
	return false
}

func memeCouleur(pixel1, pixel2 MessagePixel) bool {
	if pixel1.Rouge == pixel2.Rouge && pixel1.Vert == pixel2.Vert && pixel1.Bleu == pixel2.Bleu {
		return true
	}
	return false
}

func replaceOrAddPixel(carte []MessagePixel, newPixel MessagePixel) []MessagePixel {
	var found = false

	for i, pixel := range carte {
		if memePosition(pixel, newPixel) {
			carte[i] = newPixel
			found = true
		}
	}

	if !found {
		carte = append(carte, newPixel)
	}

	return carte
}

func getPrepostOnSamePosition(pixelReference MessagePixel, listPrepost []Message) (bool, MessagePixel) {
	var pixel MessagePixel
	var pixelFound MessagePixel
	var found = false

	for _, message := range listPrepost {
		pixel = message.Pixel
		if memePosition(pixel, pixelReference) {
			pixelFound = pixel
			found = true
			break
		}
	}

	return found, pixelFound
}
