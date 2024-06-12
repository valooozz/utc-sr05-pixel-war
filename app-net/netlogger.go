package main

import (
	"strconv"
	"utils"
)

// Mise en forme du log d'un pixel
func pixelLisible(p utils.MessagePixel) string {
	return "[(" + strconv.Itoa(p.PositionX) + "," + strconv.Itoa(p.PositionY) + ")|R" + strconv.Itoa(p.Rouge) + "|G" + strconv.Itoa(p.Vert) + "|B" + strconv.Itoa(p.Bleu) + "]"
}

// Fonction d'écriture dans les logs client
func log(str string) {
	wsSend(str)
}

// Mise en forme d'un log de l'app-net
func preparateur(typeAction string, messageNet utils.MessageNet) {
	origine := messageNet.Header.Origine
	destination := messageNet.Header.Destination
	var str string
	message := messageNet.MessageControl

	if utils.TrouverValeur(message, "couleur") != "" { //Cas d'un message de control habituel
		messageControl := utils.StringToMessage(message)
		str = "Pixel de " + messageControl.Nom[:2] + " - " + pixelLisible(messageControl.Pixel)
	} else if utils.TrouverValeur(message, "etat") != "" { //Cas d'un message etat
		messageEtat := utils.StringToMessageEtat(message)
		str = "Etat de " + messageEtat.EtatLocal.NomSite + " - (Liste des messages masquées)"
	} else if utils.TrouverValeur(message, "siteCible") != "" { //Cas d'un message accuse de reception
		messageAccuse := utils.StringToMessageAccuse(message)
		str = "Accuse de " + strconv.Itoa(messageAccuse.Estampille.Site+1) + " pour " + strconv.Itoa(messageAccuse.SiteCible+1) + " - H = " + strconv.Itoa(messageAccuse.Estampille.Horloge)
	} else if utils.TrouverValeur(message, "estampilleSite") != "" { //Cas d'un message SC
		messageSC := utils.StringToMessageExclusionMutuelle(message)
		var t string
		if messageSC.Type == utils.Requete {
			t = "Requete"
		} else {
			t = "Liberation"
		}
		str = t + " de " + strconv.Itoa(messageSC.Estampille.Site+1) + " - H = " + strconv.Itoa(messageSC.Estampille.Horloge)
	} else { //À améliorer lors de l'ajout du réseau dynamique
		str = "Message non reconnu"
	}

	log(typeAction + " : " + strconv.Itoa(origine) + "=>" + strconv.Itoa(destination) + " : " + str)
}
