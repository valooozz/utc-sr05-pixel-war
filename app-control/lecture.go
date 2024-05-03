package main

import (
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

			// Traitement des messages de contrôle
			if utils.TrouverValeur(rcvmsg, "horloge") != "" {
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" {
					traiterMessagePrepost(rcvmsg)
				} else {
					//L'affichage sur stderr se fait dans le traitement pour ce type de message
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

// TRAITEMENT DES CONTRÔLES NORMAUX : on extrait le pixel que l'on exploite dans l'app-base et on fait suivre l'information
// et tout cela avec les bonnes informations mises à jour dans le message : horloge, couleur
func traiterMessageControle(rcvmsg string) {
	monBilan--
	message := utils.StringToMessage(rcvmsg)

	if message.Nom != monNom { // On traite le message uniquement s'il ne vient pas de nous
		utils.DisplayWarning(monNom, "main", "Message de contrôle reçu : "+rcvmsg)

		// Extraction de la partie pixel
		messagePixel := message.Pixel

		// Recalage de l'horloge locale et mise à jour de sa valeur dans le message également
		H = utils.Recaler(H, message.Horloge)
		message.Horloge = H

		// Première fois qu'on reçoit l'ordre de transmettre sa sauvegarde
		if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
			maCouleur = utils.Jaune

			// On met à jour l'état local
			monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, messagePixel)

			messageEtat := utils.MessageEtat{monEtatLocal, monBilan}
			go envoyerMessage(utils.MessageEtatToString(messageEtat))

			// Réception d'un message prépost pas encore marqué comme prépost
		} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
			if jeSuisInitiateur {
				// On ajoute le message reçu à la sauvegarde générale
				etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)
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
		nbMessagesAttendus--

		// On ajoute le message prepost à la sauvegarde générale
		message := utils.StringToMessage(rcvmsg)
		etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)

		if nbEtatsAttendus == 0 && nbMessagesAttendus == 0 {
			// FIN DE L'ALGORITHME DE SAUVEGARDE
			utils.DisplayInfo(monNom, "traiterMessagePrepost", "Sauvegarde complétée")
		}
	} else {
		go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
	}
}

func traiterMessageEtat(rcvmsg string) {
	if jeSuisInitiateur {
		messageEtat := utils.StringToMessageEtat(rcvmsg)

		// On ajoute l'état local reçu à la sauvegarde générale
		etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, messageEtat.EtatLocal)

		nbEtatsAttendus--
		nbMessagesAttendus = nbMessagesAttendus - messageEtat.Bilan

		if nbEtatsAttendus == 0 && nbMessagesAttendus == 0 {
			// FIN DE L'ALGORITHME DE SAUVEGARDE
			utils.DisplayInfo(monNom, "traiterMessageEtat", "Sauvegarde complétée")
		}
	} else {
		go envoyerMessage(rcvmsg)
	}
}

func traiterMessagePixel(rcvmsg string) {
	monBilan++
	messagePixel := utils.StringToMessagePixel(rcvmsg)
	H++

	// Mise à jour de l'état local
	monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, messagePixel)

	message := utils.Message{messagePixel, H, monNom, maCouleur, false}
	go envoyerMessageControle(message)
}

func traiterDebutSauvegarde() {
	maCouleur = utils.Jaune
	jeSuisInitiateur = true
	nbEtatsAttendus = N - 1
	nbMessagesAttendus = monBilan

	// On ajoute l'état local à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, monEtatLocal)
}
