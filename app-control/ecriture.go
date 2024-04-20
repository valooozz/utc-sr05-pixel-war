package main

import (
	"fmt"
	"utils"
)

// Envoi une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	mutex.Lock()
	fmt.Println(message)
	mutex.Unlock()
}

// Envoi un type Message pour les applis de contrôle
func envoyerMessageControle(message utils.Message) {
	envoyerMessage(utils.MessageToString(message))
}

// Envoi un type MessagePixel pour l'appli de base
func envoyerMessageBase(messagePixel utils.MessagePixel) {
	envoyerMessage("A" + utils.MessagePixelToString(messagePixel))
}
