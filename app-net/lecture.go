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
			if utils.TrouverValeur(message, "header") != "" {
				traiterMessageNet(message)
			} else { //Cas de la réception d'un message de l'app-control associée
				traiterMessageControl(message)
			}
		}
		mutex.Unlock()
	}

}

func traiterMessageNet(message string) {
	utils.DisplayError(monNom, "traiterMessageNet", "Reçu : "+message)
	//à l'avenir, tous les messages ne sont pas traités de cette manière mais traités ou non
	messageNet := utils.StringToMessageNet(message)
	//header := messageNet.Header à stocker derrière
	messageControl := messageNet.MessageControl
	envoyerControl(messageControl) //envoi à l'app de control du site courant
}

func traiterMessageControl(message string) {
	header := utils.Header{ChampFictif: "contenuFictif"}
	messageNet := utils.MessageNet{Header: header, MessageControl: message}
	utils.DisplayError(monNom, "traiterMessageControl", "Envoi : "+utils.MessageNetToString(messageNet))
	envoyerNet(utils.MessageNetToString(messageNet))
}
