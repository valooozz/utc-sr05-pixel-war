package main

import (
	"fmt"
	"strconv"
	"utils"
)

func lecture() {
	var rcvmsg string

	utils.DisplayError(monNom, "Ma table & monNum ", utils.TableDeRoutageToString(tableDeRoutage)+" "+strconv.Itoa(monNum))
	for monEtat != "inactif" {
		fmt.Scanln(&rcvmsg)

		if monEtat == "inactif" {
			go envoyerMessage(rcvmsg)
			break
		}

		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu")
			break
		}

		mutex.Lock()
		//time.Sleep(time.Duration(1) * time.Second)

		utils.DisplayWarning(monNom, "lecture", "Message reçu : "+rcvmsg)

		if rcvmsg[0] == uint8('N') {
			message := rcvmsg[1:]

			if utils.TrouverValeur(message, "header") != "" {
				traiterMessageNet(message)
			} else if utils.TrouverValeur(message, "id") != "" { //Cas de la réception d'un message de l'app-control associée
				traiterMessageId(message)
			} else if utils.TrouverValeur(rcvmsg, "type") == "demande" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible == monNum {
					traiterDemandeRaccord(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "type") == "acceptation" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible == monNum && monEtat == "depart" {
					traiterAcceptationRaccord(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "type") == "depart" {
				traiterDepartRaccord()
			} else if utils.TrouverValeur(rcvmsg, "type") == "signal" {
				traiterSignalRaccord(rcvmsg)
			} else if utils.TrouverValeur(rcvmsg, "type") == "voisin" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible == monNum {
					traiterVoisinRaccord()
				}
			} else if utils.TrouverValeur(rcvmsg, "coloration") == "1" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible != monNum { // On traite un message bleu que si on n'est pas la cible (du coup le concept de cible est inversé, c'est pour gérer le fait qu'un message bleu est souvent envoyé à tous les voisins sauf un)
					traiterMessageBleu(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "coloration") == "2" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible == monNum { // On traite un message rouge que si on est la cible
					traiterMessageRouge(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "coloration") == "3" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(rcvmsg, "cible"))
				if cible != monNum { // On traite un message vert que si on n'est pas la cible
					traiterMessageVert(rcvmsg)
				}
			}
		}

		mutex.Unlock()
	}
}

func traiterMessageId(message string) {
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Reçu : "+message)
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
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Envoi : "+utils.MessageNetToString(messageNet))
	envoyerNet(utils.MessageNetToString(messageNet))
	preparateur("E", messageNet) //log au niveau du client
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Envoyé : "+utils.MessageNetToString(messageNet))
}

func traiterMessageNet(message string) {
	messageNet := utils.StringToMessageNet(message)
	header := messageNet.Header
	if header.Destination == monNum {
		preparateur("R", messageNet) //log au niveau du client
		//utils.utils.DisplayError(monNom, "traiterMessageNet", "Reçu : "+message)
		//if header.Vecteur[monNum-1] == 1 || (header.Initiateur == monNum && !utils.IlNeRestePlusQue(header.Initiateur, header.Vecteur)) || header.Origine != tableDeRoutage[0].Origine { //nième réception ou repassage par l'initiateur
		if header.Origine != tableDeRoutage[0].Origine { //nième réception ou repassage par l'initiateur
			headerForward := header
			headerForward.Destination = utils.GetDestinationFor(headerForward.Origine, tableDeRoutage)
			headerForward.Origine = monNum
			messageNet.Header = headerForward
			envoyerNet(utils.MessageNetToString(messageNet))
			preparateur("E", messageNet) //log au niveau du client
			//utils.utils.DisplayError(monNom, "traiterMessageNet", "Envoyé : "+utils.MessageNetToString(messageNet))
		} else { //Première réception : on prend en charge le message
			siteIdCpt++
			headers[strconv.Itoa(siteIdCpt)] = header
			messageControl := messageNet.MessageControl
			//Ici on vient wrapper le message dans une capsule dédiée avec un id
			messageId := utils.MessageId{Id: siteIdCpt, Message: messageControl}
			envoyerMessageId(utils.MessageIdToString(messageId)) //envoi à l'app de control du site courant
			//utils.utils.DisplayError(monNom, "traiterMessageNet", "IDENVOYÉ : "+utils.MessageIdToString(messageId))
		}
	}
}

/////////////
// Election
/////////////

func debutVague() {
	utils.DisplayWarning(monNom, "debut", "Début de la vague")
	if monParent == 0 { // Le site n’a pas encore été atteint par la vague ; il peut encore se déclarer candidat.
		monElu = monNum
		monParent = monNum

		envoyerMessageBleu(monNum) // Pour tous les voisins
	}
}

func traiterMessageBleu(rcvmsg string) {
	utils.DisplayInfoSC(monNom, "traiter", "Traitement message bleu")
	messageVague := utils.StringToMessageVague(rcvmsg)
	info := messageVague.Info
	site := messageVague.Site

	if monElu > info { // Première vague reçue, ou vague dont l’identité de l’élu est plus petite que la précédente.

		utils.DisplayInfoSC(monNom, "traiter", "Nouvelle vague reçue qui prend la place")

		monElu = info
		monParent = site
		nbVoisinsAttendus--

		if nbVoisinsAttendus > 0 {
			utils.DisplayInfoSC(monNom, "traiter", "En attente de voisins, on envoie un message bleu à tous les voisins sauf celui qui vient de nous en envoyer")
			envoyerMessageBleu(site) // Pour tous les voisins sauf j
		} else {
			utils.DisplayInfoSC(monNom, "traiter", "Tous les voisins ont répondu, on retourne un message rouge")
			envoyerMessageRouge(site) // Pour j
		}
	} else if monElu == info { // Même vague mais le site est déjà au courant
		utils.DisplayInfoSC(monNom, "traiter", "Même vague, on retourne un message rouge")
		envoyerMessageRouge(site)
	}
}

func traiterMessageRouge(rcvmsg string) {
	utils.DisplayError(monNom, "traiter", "Traitement message rouge")
	messageVague := utils.StringToMessageVague(rcvmsg)
	info := messageVague.Info

	if monElu == info { // Seuls les messages de retour appartenant à la vague en cours sont acceptés
		utils.DisplayError(monNom, "traiter", "Message de la vague en cours")
		nbVoisinsAttendus--

		if nbVoisinsAttendus == 0 {
			utils.DisplayError(monNom, "traiter", "Tous les voisins ont répondu")
			if monElu == monNum {
				traiterFinElection()
			} else {
				utils.DisplayError(monNom, "traiter", "J'envoie un message rouge à mon parent")
				envoyerMessageRouge(monParent)
			}
		}
	}
}

func traiterMessageVert(rcvmsg string) {
	utils.DisplayInfo(monNom, "traiter", "Traitement message vert")
	messageVague := utils.StringToMessageVague(rcvmsg)

	if monParent != 0 { // Si notre état n'a pas encore été réinitialisé
		reinitialiserVague(messageVague.Info)
		envoyerMessageVert(messageVague.Info, messageVague.Site)
	}
}

func traiterFinElection() {
	utils.DisplayError(monNom, "traiter", "Algo terminé, je suis élu")

	envoyerAcceptationRaccord(demande.Site)

	envoyerMessageVert(demande.Info, monNum)
	reinitialiserVague(demande.Info)
}

func reinitialiserVague(info int) {
	utils.DisplayInfo(monNom, "traiter", "Réinitialisation")
	monParent = 0
	nbVoisinsAttendus = *pVoisins
	N = N + info
	monElu = N * 100
	demande.Site = 0
	demande.Info = 0

	utils.DisplayInfoSauvegarde(monNom, "traiter", "N="+strconv.Itoa(N))
}

////////////////
// Raccordement
////////////////

func traiterDemandeRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement demande de raccord :"+rcvmsg)
	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	if demande.Site == 0 { // Si on n'a pas déjà une demande en cours
		demande.Site = messageRaccord.Site
		demande.Info = messageRaccord.Info
		debutVague()
	}
}

func traiterAcceptationRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement acceptation de raccord")

	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	if monEtat == "attente" {
		monEtat = "actif"
		N = messageRaccord.Info
		utils.DisplayInfoSauvegarde(monNom, "traiter", "N="+strconv.Itoa(N))

		go lecture()

		envoyerSignalRaccord(1, monNum)

	} else if monEtat == "depart" {
		utils.DisplayInfoSauvegarde(monNom, "traiter", "Désactivation")
		monEtat = "inactif"

		envoyerSignalRaccord(-1, monNum)

		//go transmission()
	}
}

func traiterDepartRaccord() {
	envoyerMessageVert(demande.Info, monNum)
	reinitialiserVague(demande.Info)
}

func traiterSignalRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement signal raccord")
	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	*pVoisins = *pVoisins + messageRaccord.Info
	nbVoisinsAttendus = *pVoisins

	utils.DisplayInfoSauvegarde(monNom, "traiter", "nbVoisinsAttendus="+strconv.Itoa(nbVoisinsAttendus))

	if messageRaccord.Info > 0 { // Si un site a rejoint, on lui signale notre existence
		envoyerVoisinRaccord(messageRaccord.Site)
	}
}

func traiterVoisinRaccord() {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement voisin raccord")

	*pVoisins = *pVoisins + 1
	nbVoisinsAttendus = *pVoisins

	utils.DisplayInfoSauvegarde(monNom, "traiter", "nbVoisinsAttendus="+strconv.Itoa(nbVoisinsAttendus))
}
