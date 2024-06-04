package main

import (
	"fmt"
	"utils"
)

func lecture() {
	var rcvmsg string
	for {
		fmt.Scanln(&rcvmsg)
		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu -> Fin du processus")
			continue
		}
		mutex.Lock()
		if rcvmsg[0] == uint8('N') {
			message := rcvmsg[1:]
			if utils.TrouverValeur(message, "champFictif") != "" {
				traiterMessageNet(message)
			} else { //Cas de la réception d'un message de l'app-control associée
				traiterMessageControl(message)
			}
		}
		mutex.Unlock()
	}

}

func traiterMessageNet(message string) {
	//à l'avenir, tous les messages ne sont pas traités de cette manière mais traités ou non
	messageNet := utils.StringToMessageNet(message)
	messageControl := messageNet.MessageControl
	envoyerControl(messageControl) //envoi à l'app de control du site courant
}

func traiterMessageControl(message string) {
	messageNet := utils.MessageNet{ChampFictif: "contenuFictif", MessageControl: message}
	envoyerNet(utils.MessageNetToString(messageNet))
}
