package main

import (
	"fmt"
	"strconv"
	"utils"
)

// Fonction de lecture d'une application inactive qui n'a pas encore rejoint le réseau
func attenteRaccordement() {
	var rcvmsg string

	for monEtat != "actif" {
		fmt.Scanln(&rcvmsg) // On reçoit les messages même si on est inactifs, pour les évacuer de notre entrée standard

		if rcvmsg[0] == uint8('C') { // On ignore les messages provenant d'une app de contrôle
			continue
		}

		message := rcvmsg[1:] // On extrait le message en retirant la lettre préfixe

		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu")
			break
		}

		mutex.Lock()

		// Si on reçoit une acceptation de raccord, on la traite
		if utils.TrouverValeur(message, "type") == "acceptation" {
			cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
			if cible == monNum {
				traiterAcceptationRaccord(message)
			}
		}

		mutex.Unlock()
	}
}

// Fonction de lecture d'une application active dans le réseau
func lecture() {
	var rcvmsg string

	utils.DisplayError(monNom, "Ma table & monNum ", utils.TableDeRoutageToString(tableDeRoutage)+" "+strconv.Itoa(monNum))

	for monEtat != "inactif" {
		fmt.Scanln(&rcvmsg)

		// On peut passer inactif lors de l'attente d'un message
		if monEtat == "inactif" {
			break
		}

		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu")
			break
		}

		mutex.Lock()

		//utils.DisplayWarning(monNom, "lecture", "Message reçu : "+rcvmsg)

		// On ne traite que les messages portant le préfixe 'N'
		if rcvmsg[0] == uint8('N') {
			message := rcvmsg[1:] // On extrait le message en retirant le préfixe 'N'

			if utils.TrouverValeur(message, "header") != "" { // Cas d'un message d'une app-net qui transmet des informations sur l'anneau logique
				traiterMessageNet(message)
			} else if utils.TrouverValeur(message, "id") != "" { // Cas d'un message de l'app-control associée
				traiterMessageId(message)
			} else if utils.TrouverValeur(message, "type") == "demande" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible == monNum {
					traiterDemandeRaccord(message)
				}
			} else if utils.TrouverValeur(message, "type") == "acceptation" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible == monNum && monEtat == "depart" {
					traiterAcceptationRaccord(message)
				}
			} else if utils.TrouverValeur(message, "type") == "depart" {
				traiterDepartRaccord()
			} else if utils.TrouverValeur(message, "type") == "signal" {
				traiterSignalRaccord(message)
			} else if utils.TrouverValeur(message, "type") == "voisin" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible == monNum {
					traiterVoisinRaccord()
				}
			} else if utils.TrouverValeur(message, "coloration") == "1" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible != monNum { // On traite un message bleu que si on n'est pas la cible (du coup le concept de cible est inversé, c'est pour gérer le fait qu'un message bleu est souvent envoyé à tous les voisins sauf un)
					traiterMessageBleu(message)
				}
			} else if utils.TrouverValeur(message, "coloration") == "2" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible == monNum { // On traite un message rouge que si on est la cible
					traiterMessageRouge(message)
				}
			} else if utils.TrouverValeur(message, "coloration") == "3" {
				cible, _ := strconv.Atoi(utils.TrouverValeur(message, "cible"))
				if cible != monNum { // On traite un message vert que si on n'est pas la cible
					traiterMessageVert(message)
				}
			}
		}

		mutex.Unlock()
	}

	// Si on arrive là, c'est qu'on est passé inactif, on transmet donc le premier message reçu
	transmission(rcvmsg)
}

// Fonction de lecture pour un site inactif dans le réseau, qui doit retransmettre les messages qu'il reçoit
func transmission(rcvmsg string) {
	fmt.Println(rcvmsg)

	for monEtat == "inactif" {
		fmt.Scanln(&rcvmsg)

		// On ne traite que les messages portant le préfixe 'N'
		if rcvmsg[0] == uint8('N') {
			message := rcvmsg[1:] // On extrait le message en retirant le préfixe 'N'

			// Si c'est un message pour l'anneau logique, on le transmet sur l'anneau
			if utils.TrouverValeur(message, "header") != "" {
				messageNet := utils.StringToMessageNet(message)
				header := messageNet.Header
				if header.Destination == monNum { // Si le message nous est destiné sur l'anneau
					headerForward := header
					headerForward.Destination = utils.GetDestinationFor(headerForward.Origine, tableDeRoutage)
					headerForward.Origine = monNum
					messageNet.Header = headerForward
					envoyerNet(utils.MessageNetToString(messageNet))
					preparateur("E", messageNet) //log au niveau du client
				}
				// Si c'est un message pour les apps net, et qu'on a plus (+) d'un voisin (évite quelques cas de ping-pong infini)
			} else if nbVoisinsAttendus > 1 {
				fmt.Println(rcvmsg)
			}
		}
	}
}

//////////////////////////////////
// Diffusion sur l'anneau logique
//////////////////////////////////

// Traite un message reçu de l'app-control, qui porte donc un id
func traiterMessageId(message string) {
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Reçu : "+message)
	messageId := utils.StringToMessageId(message)
	if messageId.Message == "" { //Cas d'un message arreté par l'application de control et renvoyé avec un id mais un message vide
		delete(headers, strconv.Itoa(messageId.Id))
		return
	}
	var header utils.Header
	if messageId.Id == -1 {
		// Il faut le wrapper pour la première fois
		vecteur := make([]int, N)
		// Initialisation du header avec num du site courant, destination de la première règle de routage en destination, etc.
		header = utils.Header{Origine: monNum, Destination: tableDeRoutage[0].Destination, Initiateur: monNum, Vecteur: vecteur}
	} else {
		header = headers[strconv.Itoa(messageId.Id)]
		delete(headers, strconv.Itoa(messageId.Id))
		//header.Vecteur[monNum-1] = 1
		header.Destination = utils.GetDestinationFor(header.Origine, tableDeRoutage)
		header.Origine = monNum
		// Il faut récupérer son header dans la map headers pour le wrapper avec (ne pas oublier de maj certains headers)
	}
	messageNet := utils.MessageNet{Header: header, MessageControl: messageId.Message}
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Envoi : "+utils.MessageNetToString(messageNet))
	envoyerNet(utils.MessageNetToString(messageNet))
	preparateur("E", messageNet) //log au niveau du client
	//utils.utils.DisplayError(monNom, "traiterMessageId", "Envoyé : "+utils.MessageNetToString(messageNet))
}

// Traite un message reçu d'une app-net
func traiterMessageNet(message string) {
	messageNet := utils.StringToMessageNet(message)
	header := messageNet.Header

	if header.Destination == monNum { // Si le message nous est destiné sur l'anneau
		preparateur("R", messageNet) // log au niveau du client
		//utils.utils.DisplayError(monNom, "traiterMessageNet", "Reçu : "+message)
		if header.Origine != tableDeRoutage[0].Origine { // nième réception ou repassage par l'initiateur
			headerForward := header
			headerForward.Destination = utils.GetDestinationFor(headerForward.Origine, tableDeRoutage)
			headerForward.Origine = monNum
			messageNet.Header = headerForward
			envoyerNet(utils.MessageNetToString(messageNet))
			preparateur("E", messageNet) //log au niveau du client
			//utils.utils.DisplayError(monNom, "traiterMessageNet", "Envoyé : "+utils.MessageNetToString(messageNet))
		} else { // Première réception : on prend en charge le message
			siteIdCpt++
			headers[strconv.Itoa(siteIdCpt)] = header
			messageControl := messageNet.MessageControl
			//Ici on vient wrapper le message dans une capsule dédiée avec un id
			messageId := utils.MessageId{Id: siteIdCpt, Message: messageControl}
			envoyerMessageId(utils.MessageIdToString(messageId)) // On transmet le message à son app-control
			//utils.utils.DisplayError(monNom, "traiterMessageNet", "IDENVOYÉ : "+utils.MessageIdToString(messageId))
		}
	}
}

/////////////
// Election
/////////////

// Lance une vague (pour l'algorithme d'élection par extinction de vague)
func debutVague() {
	utils.DisplayWarning(monNom, "debut", "Début de la vague")
	if monParent == 0 { // Le site n’a pas encore été atteint par la vague, il peut encore se déclarer candidat
		monElu = monNum
		monParent = monNum

		envoyerMessageBleu(monNum) // Pour tous les voisins
	}
}

// Traite un message bleu (message descendant l'arborescence)
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
			envoyerMessageBleu(site) // Pour tous les voisins sauf le nouveau parent
		} else {
			utils.DisplayInfoSC(monNom, "traiter", "Tous les voisins ont répondu, on retourne un message rouge")
			envoyerMessageRouge(site) // Pour le nouveau parent
		}
	} else if monElu == info { // Même vague mais le site est déjà au courant
		utils.DisplayInfoSC(monNom, "traiter", "Même vague, on retourne un message rouge")
		envoyerMessageRouge(site) // Pour le site qui a envoyé le message bleu
	}
}

// Traite un message rouge (message remontant l'arborescence)
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

// Traite un message vert, qui correspond à un demi-vague servant à réinitialiser le réseau pour la prochaine élection
func traiterMessageVert(rcvmsg string) {
	utils.DisplayInfo(monNom, "Traitement message vert", rcvmsg)
	messageVague := utils.StringToMessageVague(rcvmsg)

	if monParent != 0 { // Si notre état n'a pas encore été réinitialisé
		reinitialiserVague(messageVague.Info, messageVague.SiteDemandeur)
		envoyerMessageVert(messageVague.Info, messageVague.Site, messageVague.SiteDemandeur)
	}
}

// Traite la fin d'une élection, quand on est élu
func traiterFinElection() {
	utils.DisplayError(monNom, "traiter", "Algo terminé, je suis élu")

	envoyerAcceptationRaccord(demande.Site) // On envoie une acceptation au site demandeur

	if demande.Info == 1 { // Pour le moment, on ne change le routage du site élu qu'à l'arrivée d'un membre, au même titre que les autres
		majRoutage(demande.Site)
	}

	envoyerMessageVert(demande.Info, monNum, demande.Site)
	reinitialiserVague(demande.Info, demande.Site)
}

// Met à jour la table de routage pour intégrer le nouveau site
func majRoutage(nouveauSite int) {
	tmp := tableDeRoutage[0].Destination
	tableDeRoutage[0].Destination = nouveauSite
	tableDeRoutage = append(tableDeRoutage, utils.Route{nouveauSite, tmp})
	utils.DisplayError(monNom, "majRoutage", utils.TableDeRoutageToString(tableDeRoutage))
}

// Réinitialise les informations associées à l'algorithme d'élection par extinction de vague
func reinitialiserVague(info int, siteDemandeur int) {
	utils.DisplayInfo(monNom, "traiter", "Réinitialisation")
	monParent = 0
	nbVoisinsAttendus = *pVoisins
	N = N + info
	//utils.DisplayError(monNom, "Je remonterais ", strconv.Itoa(siteDemandeur))

	// On envoie le nouveau N à l'app-control, ainsi que le numéro du site qui rejoint ou part
	envoyerSpecialControl("N=" + strconv.Itoa(N) + "|" + strconv.Itoa(siteDemandeur))

	monElu = N * 100
	demande.Site = 0
	demande.Info = 0

	utils.DisplayInfoSauvegarde(monNom, "traiter", "N="+strconv.Itoa(N))
}

////////////////
// Raccordement
////////////////

// Traite une demande de raccord
func traiterDemandeRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement demande de raccord :"+rcvmsg)
	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	if demande.Site == 0 { // Si on n'a pas déjà une demande en cours
		demande.Site = messageRaccord.Site
		demande.Info = messageRaccord.Info
		debutVague()
	}
}

// Traite une acceptation de raccord
func traiterAcceptationRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement acceptation de raccord")

	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	if monEtat == "attente" { // Si on rejoint
		monEtat = "actif"
		N = messageRaccord.Info
		utils.DisplayInfoSauvegarde(monNom, "traiter", "N="+strconv.Itoa(N))
		envoyerSpecialControl("NN=" + strconv.Itoa(N)) // On transmet à l'app-control le nombre de sites sur le réseau

		go lecture()

		envoyerSignalRaccord(1, monNum) // On prévient les voisins de notre arrivée sur le réseau

	} else if monEtat == "depart" { // Si on part
		utils.DisplayInfoSauvegarde(monNom, "traiter", "Désactivation")
		monEtat = "inactif"

		envoyerSignalRaccord(-1, monNum) // On prévient les voisins de notre départ du réseau
	}
}

func traiterDepartRaccord() {
	envoyerMessageVert(demande.Info, monNum, demande.Site)
	reinitialiserVague(demande.Info, demande.Site)
}

// Traite un signal envoyé par un site qui rejoint ou part
func traiterSignalRaccord(rcvmsg string) {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement signal raccord")
	messageRaccord := utils.StringToMessageRaccord(rcvmsg)

	// On met à jour le nombre de voisins
	*pVoisins = *pVoisins + messageRaccord.Info
	nbVoisinsAttendus = *pVoisins

	utils.DisplayInfoSauvegarde(monNom, "traiter", "nbVoisinsAttendus="+strconv.Itoa(nbVoisinsAttendus))

	if messageRaccord.Info > 0 { // Si un site a rejoint, on lui signale notre existence
		envoyerVoisinRaccord(messageRaccord.Site)
	}
}

// Traite un signal envoyé par les voisins après qu'on les a prévenus de notre arrivée
func traiterVoisinRaccord() {
	utils.DisplayInfoSauvegarde(monNom, "traiter", "Traitement voisin raccord")

	// On met à jour notre nombre de voisins
	*pVoisins = *pVoisins + 1
	nbVoisinsAttendus = *pVoisins

	utils.DisplayInfoSauvegarde(monNom, "traiter", "nbVoisinsAttendus="+strconv.Itoa(nbVoisinsAttendus))
}
