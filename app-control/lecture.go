package main

import (
	"container/list"
	"fmt"
	"utils"
)

// Pour l'instant, boucle sur l'entrée standard, lit et communique le résultat à la routine d'écriture
func lecture() {
	var rcvmsg string
	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		// On traite uniquement les messages qui ne commencent pas par un 'A'
		if rcvmsg[0] != uint8('A') {
			//TRAITEMENT DES MESSAGES DE CONTRÔLE
			if utils.TrouverValeur(rcvmsg, "horloge") != "" {
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" {
					traiterMessagePrepost(rcvmsg)
				} else {
					utils.DisplayWarning(monNom, "main", "Message de contrôle reçu : "+rcvmsg)
					traiterMessageControle(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "etat") != "" {
				traiterMessageEtat(rcvmsg)
			} else {
				utils.DisplayWarning(monNom, "lecture", "Message pixel reçu : "+rcvmsg)
				traiterMessagePixel(rcvmsg)
			}
		}
		mutex.Unlock()
	}
}

//TRAITEMENT DES CONTRÔLES NORMAUX : on extrait le pixel que l'on exploite dans l'app-base et on fait suivre l'information
// et tout cela avec les bonnes informations mises à jour dans le message : horloge, couleur
func traiterMessageControle(rcvmsg string) {
	message := utils.StringToMessage(rcvmsg)

	if message.Nom != monNom { // On traite le message uniquement s'il ne vient pas de nous
		//Extraction de la partie pixel
		messagePixel := message.Pixel
		//Recalage de l'horloge locale et mise à jour de sa valeur dans le message également
		H = utils.Recaler(H, message.Horloge)
		message.Horloge = H

		//ATTENTION ICI, METTRE À JOUR L'ÉTAT GLOBAL AVANT D'ENVOYER QUOI QUE CE SOIT

		//Avertissement d'une coupure demandée et actions en conséquence
		if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
			maCouleur = utils.Jaune
			messageEtat := utils.MessageEtat{list.List(monEtatLocal), monBilan}
			go envoyerMessage(utils.MessageEtatToString(messageEtat))
			//Réception d'un message prépost pas encore marqué comme prépost
		} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
			if jeSuisInitiateur {
				// Ajouter le message reçu à la sauvegarde générale
			} else {
				messagePrepost := message
				messagePrepost.Prepost = true
				go envoyerMessageControle(messagePrepost)
			}
		}

		message.Couleur = maCouleur
		go envoyerMessageControle(message)  // Pour la prochaine app de contrôle de l'anneau
		go envoyerMessageBase(messagePixel) // Pour l'app de base
	}
}

func traiterMessagePrepost(rcvmsg string) {
	if jeSuisInitiateur {
		//message := utils.StringToMessage(rcvmsg)
		// Traiter l'ajout du message à l'état de sauvegarde
	} else {
		go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
	}
}

func traiterMessageEtat(rcvmsg string) {
	if jeSuisInitiateur {
		// Traiter l'ajout de l'état à la sauvegarde générale
	} else {
		go envoyerMessage(rcvmsg)
	}
}

func traiterMessagePixel(rcvmsg string) {
	messagePixel := utils.StringToMessagePixel(rcvmsg)
	H++
	//ATTENTION ICI, METTRE À JOUR L'ÉTAT GLOBAL
	message := utils.Message{messagePixel, H, monNom, maCouleur, false}
	go envoyerMessageControle(message)
}