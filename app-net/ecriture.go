package main

import (
	"fmt"
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
