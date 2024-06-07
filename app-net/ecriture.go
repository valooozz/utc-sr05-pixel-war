package main

import (
	"fmt"
	"time"
	"utils"
)

// Envoi une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	//mutex.Lock()
	fmt.Println(message)
	//mutex.Unlock()
}

func envoyerMessageId(message string) {
	msg := "C" + message
	envoyerMessage(msg)
}

func envoyerNet(message string) {
	msg := "N" + message
	envoyerMessage(msg)
}

func envoyerSpecialControl(message string) {
	msg := "C" + message
	envoyerMessage(msg)
}

//////////////
// Election
//////////////

func envoyerMessageBleu(cible int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(1), monElu, cible}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

func envoyerMessageRouge(cible int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(2), monElu, cible}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

func envoyerMessageVert(info int, cible int) {
	messageVague := utils.MessageVague{monNum, utils.ColorationVague(3), info, cible}
	str := utils.MessageVagueToString(messageVague)

	envoyerNet(str)
}

////////////////
// Raccordement
////////////////

func envoyerDemandeRaccord(info int, cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "demande", info, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	for monEtat == "attente" || monEtat == "depart" {
		envoyerNet(str)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func envoyerAcceptationRaccord(cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "acceptation", N + 1, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}

func envoyerSignalRaccord(info int, cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "signal", info, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}

func envoyerVoisinRaccord(cible int) {
	messageRaccord := utils.MessageRaccord{monNum, "voisin", 0, cible}
	str := utils.MessageRaccordToString(messageRaccord)

	envoyerNet(str)
}
