package utils

import (
	"math"
	"strconv"
	"strings"
)

// Trouve une valeur dans un message string transmis sur l'anneau
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

// Recale une horloge entière
func Recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}

// Met à jour l'horloge vectorielle locale avec celle reçue
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

// Retourne une copie d'une horloge vectorielle (utile à cause du fonctionnement des slices en go)
func CopyHorlogeVectorielle(horlogeVectorielle HorlogeVectorielle) HorlogeVectorielle {

	var copie = HorlogeVectorielle{}

	for key, val := range horlogeVectorielle {
		copie[key] = val
	}

	return copie
}

// Retourne Vrai si la coupure présente dans l'état global est cohérente, faux sinon
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

// Met à jour l'état local en ajoutant ou remplaçant un MessagePixel
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

// Retourne une copie de l'état local entré (utile à cause du fonctionnement des slices en go)
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

// Retourne une grille de pixels unique à partir d'un état global
func ReconstituerCarte(etatGlobal EtatGlobal) []MessagePixel {
	var carte = etatGlobal.ListEtatLocal[0].ListMessagePixel
	var pixel MessagePixel

	for _, message := range etatGlobal.ListMessagePrepost {
		pixel = message.Pixel
		carte = replaceOrAddPixel(carte, pixel)
	}

	return carte
}

/////////////////////
// Exclusion mutuelle
/////////////////////

func QuestionEntreeSC(site int, tabSC []MessageExclusionMutuelle) bool {
	cpt := 0

	if tabSC[site].Type != Requete {
		return false
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
		return true
	}
	return false
}

func InitialisationNumSite(site string) int {
	StartNumberIndex := 1
	SiteString := site[StartNumberIndex:len(site)]
	NumSite, _ := strconv.Atoi(SiteString)
	return NumSite
}

func PlusVieilleRequeteAlive(site int, tabSC []MessageExclusionMutuelle) int {
	sitePrioritaire := Estampille{Site: math.MaxInt, Horloge: math.MaxInt}
	for numOtherSite := 0; numOtherSite < len(tabSC); numOtherSite++ {
		if numOtherSite == site || tabSC[numOtherSite].Type != Requete {
			continue
		}
		if tabSC[numOtherSite].Estampille.Horloge < sitePrioritaire.Horloge {
			sitePrioritaire = tabSC[numOtherSite].Estampille
		} else if tabSC[numOtherSite].Estampille.Horloge == sitePrioritaire.Horloge && tabSC[numOtherSite].Estampille.Site < sitePrioritaire.Site {
			sitePrioritaire = tabSC[numOtherSite].Estampille
		}
	}
	return sitePrioritaire.Site
}

////////////////////
// FONCTIONS PRIVEES
////////////////////

// Retourne Vrai si les deux pixels entrés sont à la même position, Faux sinon
func memePosition(pixel1, pixel2 MessagePixel) bool {
	if pixel1.PositionX == pixel2.PositionX && pixel1.PositionY == pixel2.PositionY {
		return true
	}
	return false
}

// Retourne Vrai si les deux pixels entrés sont de la même couleur, Faux sinon
func memeCouleur(pixel1, pixel2 MessagePixel) bool {
	if pixel1.Rouge == pixel2.Rouge && pixel1.Vert == pixel2.Vert && pixel1.Bleu == pixel2.Bleu {
		return true
	}
	return false
}

// Met à jour une liste de pixels en remplaçant le pixel déjà présent à la même position, ou en l'ajoutant si la position n'y est pas
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

func GetDestinationFor(source int, routage TableDeRoutage) int {
	for _, route := range routage {
		if route.Origine == source {
			return route.Destination
		}
	}
	return -1
}

func IlNeRestePlusQue(init int, vect []int) bool {
	if vect[init-1] == 1 {
		DisplayError("Utils", "IlNeRestePlusQue()", "L'initiateur n'est même pas faux")
		return false
	}
	for i, v := range vect {
		if i != init-1 { //On ne traite que les cases qui ne concernent pas l'initiateur car l'initiateur est traité en haut
			if v != 1 {
				return false
			}
		}
	}
	return true
}
