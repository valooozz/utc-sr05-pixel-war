package main

import (
	"fmt"
	"strconv"
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
			} else if utils.TrouverValeur(message, "id") != "" { //Cas de la réception d'un message de l'app-control associée
				traiterMessageId(message)
			}
		}
		mutex.Unlock()
	}

}

func traiterMessageId(message string) {
	//utils.DisplayError(monNom, "traiterMessageId", "Reçu : "+message)
	messageId := utils.StringToMessageId(message)
	if messageId.Message == "" { //Cas d'un message arreté par l'application de control et renvoyé avec un id mais un message vide
		delete(headers, strconv.Itoa(messageId.Id))
		return
	}
	var header utils.Header
	if messageId.Id == -1 {
		//Il faut le wrapper pour la première fois
		header = utils.Header{ChampFictif: "contenuFictif"}
	} else {
		header = headers[strconv.Itoa(siteIdCpt)]
		//Il faut récupérer son header dans la map headers pour le wrapper avec (ne pas oublier de maj certains headers)
		//MAJ AUTRES
	}
	messageNet := utils.MessageNet{Header: header, MessageControl: messageId.Message}
	//utils.DisplayError(monNom, "traiterMessageId", "Envoi : "+utils.MessageNetToString(messageNet))
	envoyerNet(utils.MessageNetToString(messageNet))
}

func traiterMessageNet(message string) {
	//utils.DisplayError(monNom, "traiterMessageNet", "Reçu : "+message)
	//à l'avenir, tous les messages ne sont pas traités de cette manière mais traités ou non
	messageNet := utils.StringToMessageNet(message)
	header := messageNet.Header
	siteIdCpt++
	headers[strconv.Itoa(siteIdCpt)] = header
	messageControl := messageNet.MessageControl
	//Ici on vient wrapper le message dans une capsule dédiée avec un id
	messageId := utils.MessageId{Id: siteIdCpt, Message: messageControl}
	envoyerMessageId(utils.MessageIdToString(messageId)) //envoi à l'app de control du site courant
}
