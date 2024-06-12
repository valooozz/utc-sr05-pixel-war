package main

import (
	"fmt"
	"time"
	"utils"
)

// Envoie une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	//mutex.Lock()
	fmt.Println(message)
	//mutex.Unlock()
}

// Ajoute un 'C' au début du message pour qu'il soit traité par l'app de contrôle
func envoyerMessageId(message string) {
	msg := "C" + message
	envoyerMessage(msg)
}

// Ajoute un 'N' au début du message pour qu'il soit traité par les autres app network
func envoyerNet(message string) {
	msg := "N" + message
	envoyerMessage(msg)
}

// Ajoute un 'C' au début du message pour qu'il soit traité par l'app de contrôle
func envoyerSpecialControl(message string) {
	msg := "C" + message
	envoyerMessage(msg)
}

////////////
// Election
////////////

// Envoie un message bleu avec la cible spécifiée
func envoyerMessageBleu(cible int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(1), monElu, cible, -1}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

// Envoie un message rouge avec la cible spécifiée
func envoyerMessageRouge(cible int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(2), monElu, cible, -1}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

// Envoie un message vert avec le sens (+1 ou -1) de raccordement, la cible, et le site demandeur du raccordement
func envoyerMessageVert(info int, cible int, siteDemandeur int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(3), info, cible, siteDemandeur}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

////////////////
// Raccordement
////////////////

// Envoie une demande de raccordement, avec le sens (+1 pour rejoindre et -1 pour partir) et la cible à qui faire la demande
func envoyerDemandeRaccord(info int, cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "demande", info, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	for monEtat == "attente" || monEtat == "depart" {
		envoyerNet(str)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

// Envoie une acceptation du raccord à la cible spécifiée
func envoyerAcceptationRaccord(cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "acceptation", N + 1, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}

// Envoie un signal aux voisins pour dire si on vient de rejoindre (info=1) ou de partir (info=-1)
func envoyerSignalRaccord(info int, cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "signal", info, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}

// Envoie un message pour signaler sa présence en tant que voisin au site qui vient de rejoindre
func envoyerVoisinRaccord(cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "voisin", 0, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}
