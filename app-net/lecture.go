package main

import (
	"fmt"
	"strconv"
	"utils"
)

func lecture() {
	var rcvmsg string
	utils.DisplayError(monNom, "Ma table & monNum ", utils.TableDeRoutageToString(tableDeRoutage)+" "+strconv.Itoa(monNum))
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
		vecteur := make([]int, N)
		//Initialisation du header avec num du site courant, destination de la première règle de routage en destination, etc.
		header = utils.Header{Origine: monNum, Destination: tableDeRoutage[0].Destination, Initiateur: monNum, Vecteur: vecteur}
	} else {
		header = headers[strconv.Itoa(messageId.Id)]
		delete(headers, strconv.Itoa(messageId.Id))
		header.Vecteur[monNum-1] = 1
		header.Destination = utils.GetDestinationFor(header.Origine, tableDeRoutage)
		header.Origine = monNum
		//Il faut récupérer son header dans la map headers pour le wrapper avec (ne pas oublier de maj certains headers)
	}
	messageNet := utils.MessageNet{Header: header, MessageControl: messageId.Message}
	//utils.DisplayError(monNom, "traiterMessageId", "Envoi : "+utils.MessageNetToString(messageNet))
	envoyerNet(utils.MessageNetToString(messageNet))
	preparateur("E", messageNet) //log au niveau du client
	//utils.DisplayError(monNom, "traiterMessageId", "Envoyé : "+utils.MessageNetToString(messageNet))
}

func traiterMessageNet(message string) {
	messageNet := utils.StringToMessageNet(message)
	header := messageNet.Header
	if header.Destination == monNum {
		preparateur("R", messageNet) //log au niveau du client
		//utils.DisplayError(monNom, "traiterMessageNet", "Reçu : "+message)
		//if header.Vecteur[monNum-1] == 1 || (header.Initiateur == monNum && !utils.IlNeRestePlusQue(header.Initiateur, header.Vecteur)) || header.Origine != tableDeRoutage[0].Origine { //nième réception ou repassage par l'initiateur
		if header.Origine != tableDeRoutage[0].Origine { //nième réception ou repassage par l'initiateur
			headerForward := header
			headerForward.Destination = utils.GetDestinationFor(headerForward.Origine, tableDeRoutage)
			headerForward.Origine = monNum
			messageNet.Header = headerForward
			envoyerNet(utils.MessageNetToString(messageNet))
			preparateur("E", messageNet) //log au niveau du client
			//utils.DisplayError(monNom, "traiterMessageNet", "Envoyé : "+utils.MessageNetToString(messageNet))
		} else { //Première réception : on prend en charge le message
			siteIdCpt++
			headers[strconv.Itoa(siteIdCpt)] = header
			messageControl := messageNet.MessageControl
			//Ici on vient wrapper le message dans une capsule dédiée avec un id
			messageId := utils.MessageId{Id: siteIdCpt, Message: messageControl}
			envoyerMessageId(utils.MessageIdToString(messageId)) //envoi à l'app de control du site courant
			//utils.DisplayError(monNom, "traiterMessageNet", "IDENVOYÉ : "+utils.MessageIdToString(messageId))
		}
	}
}
