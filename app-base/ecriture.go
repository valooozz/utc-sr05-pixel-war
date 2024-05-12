package main

import (
	"fmt"
	"time"
	"utils"
)

// Envoi une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	mutex.Lock()
	fmt.Println(message)
	mutex.Unlock()
}

func envoyerPixel(positionX int, positionY int, rouge int, vert int, bleu int) {
	messagePixel := utils.MessagePixel{positionX, positionY, rouge, vert, bleu}
	envoyerMessage(utils.MessagePixelToString(messagePixel))
}

// Communication APP BASE -> APP CONTROLE
func demandeSC() {
	msg := utils.MessageTypeSCToString(utils.Requete)
	envoyerMessage(msg)

	for accesSC == false {
		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func relacherSC() {
	accesSC = false
	msg := utils.MessageTypeSCToString(utils.Liberation)
	envoyerMessage(msg)
}
