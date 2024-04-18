package main

import "fmt"

// Pour l'instant, boucle sur le channel lié à la lecture puis écrit sur la sortie standard ou autre part
func envoyer_message(message string) {
	mutex.Lock()
	fmt.Print(message)
	mutex.Unlock()
}
