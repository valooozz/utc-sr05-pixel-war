package main

import (
	"fmt"
	"strconv"
	"utils"
)

// Pour l'instant, boucle sur l'entrée standard, lit et communique le résultat à la routine d'écriture
func lecture() {
	var rcvmsg string
	for {
		fmt.Scanln(&rcvmsg)
		if rcvmsg == "" {
			utils.DisplayWarning(monNom, "lecture", "Message vide reçu")
			continue
		}
		mutex.Lock()
		// On traite uniquement les messages qui ne commencent pas par un 'A'
		if rcvmsg[0] != uint8('A') {

			// Demande de sauvegarde
			if rcvmsg == "sauvegarde" {
				traiterDebutSauvegarde()

				// Traitement des messages de contrôle
			} else if utils.TrouverValeur(rcvmsg, "horloge") != "" {
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" {
					traiterMessagePrepost(rcvmsg)
				} else {
					//L'affichage sur stderr se fait dans le traitement pour ce type de message
					traiterMessageControle(rcvmsg)
				}
			} else if utils.TrouverValeur(rcvmsg, "etat") != "" {
				traiterMessageEtat(rcvmsg)
			} else {
				traiterMessagePixel(rcvmsg)
			}
		}
		mutex.Unlock()
	}
}

// TRAITEMENT DES CONTRÔLES NORMAUX : on extrait le pixel que l'on exploite dans l'app-base et on fait suivre l'information
// et tout cela avec les bonnes informations mises à jour dans le message : horloge, couleur
func traiterMessageControle(rcvmsg string) {
	message := utils.StringToMessage(rcvmsg)

	// On traite le message uniquement s'il ne vient pas de nous
	if message.Nom == monNom {
		return
	}
	monBilan--

	utils.DisplayWarning(monNom, "Controle", "Message de contrôle reçu : "+rcvmsg+" monBilanActuel = "+strconv.Itoa(monBilan))

	// Extraction de la partie pixel
	messagePixel := message.Pixel

	// Recalage de l'horloge locale et mise à jour de sa valeur dans le message également
	H = utils.Recaler(H, message.Horloge)
	message.Horloge = H

	// Mise à jour de l'horloge vectorielle locale et mise à jour de sa valeur dans le message également
	horlogeVectorielle = utils.MajHorlogeVectorielle(monNom, horlogeVectorielle, message.Vectorielle)
	message.Vectorielle = horlogeVectorielle

	// Première fois qu'on reçoit l'ordre de transmettre sa sauvegarde
	if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
		maCouleur = utils.Jaune
		utils.DisplayError(monNom, "Controle", "Passage en jaune")

		messageEtat := utils.MessageEtat{monEtatLocal, monBilan}
		utils.DisplayError(monNom, "Controle", "Etat : "+utils.MessageEtatToString(messageEtat))
		go envoyerMessageEtat(messageEtat)

		// Réception d'un message prépost pas encore marqué comme prépost
	} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
		if jeSuisInitiateur {
			// On ajoute le message reçu à la sauvegarde générale
			etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)
		} else {
			utils.DisplayError(monNom, "Controle", "Prepost")
			messagePrepost := message
			messagePrepost.Prepost = true
			go envoyerMessageControle(messagePrepost)
			monBilan++
		}
	}

	message.Couleur = maCouleur

	// On met à jour l'état local
	monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, messagePixel)
	monEtatLocal.Vectorielle = horlogeVectorielle

	go envoyerMessageControle(message) // Pour la prochaine app de contrôle de l'anneau
	monBilan++
	go envoyerMessageBase(messagePixel) // Pour l'app de base
	utils.DisplayInfo(monNom, "Controle", "monBilanActuel = "+strconv.Itoa(int(monBilan)))
}

func traiterMessagePrepost(rcvmsg string) {

	if !jeSuisInitiateur {
		go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
	}

	nbMessagesAttendus--

	// On ajoute le message prepost à la sauvegarde générale
	message := utils.StringToMessage(rcvmsg)
	etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)

	if nbEtatsAttendus == 0 && nbMessagesAttendus == 0 {
		finSauvegarde()
	}
}

func traiterMessageEtat(rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayError(monNom, "Etat", "Transfert message etat : "+rcvmsg)
		go envoyerMessage(rcvmsg)
		return
	}

	utils.DisplayError(monNom, "Etat", "MessageEtat recu")
	messageEtat := utils.StringToMessageEtat(rcvmsg)

	// On ajoute l'état local reçu à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, messageEtat.EtatLocal)

	nbEtatsAttendus--
	nbMessagesAttendus = nbMessagesAttendus + messageEtat.Bilan

	utils.DisplayError(monNom, "Etat", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus)+" ; nbMessagesAttendus="+strconv.Itoa(nbMessagesAttendus))
	if nbEtatsAttendus == 0 && nbMessagesAttendus == 0 {
		finSauvegarde()
	}
}

func traiterMessagePixel(rcvmsg string) {
	utils.DisplayWarning(monNom, "lecture", "Message pixel reçu : "+rcvmsg)

	messagePixel := utils.StringToMessagePixel(rcvmsg)
	H++

	horlogeVectorielle[monNom]++

	// Mise à jour de l'état local
	monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, messagePixel)
	monEtatLocal.Vectorielle = horlogeVectorielle

	message := utils.Message{messagePixel, H, horlogeVectorielle, monNom, maCouleur, false}
	go envoyerMessageControle(message)
	monBilan++
}

func traiterDebutSauvegarde() {
	utils.DisplayError(monNom, "Debut", "debut de la sauvegarde")
	maCouleur = utils.Jaune
	jeSuisInitiateur = true
	nbEtatsAttendus = N - 1
	nbMessagesAttendus = monBilan

	utils.DisplayError(monNom, "Debut", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus)+" ; nbMessagesAttendus="+strconv.Itoa(nbMessagesAttendus))

	// On ajoute l'état local à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, monEtatLocal)
}

func finSauvegarde() {
	utils.DisplayError(monNom, "Fin", "Sauvegarde complétée")
	for _, etatLocal := range etatGlobal.ListEtatLocal {
		utils.DisplayInfo(monNom, "Fin", utils.EtatLocalToString(etatLocal))
	}

	if utils.CoupureEstCoherente(etatGlobal) {
		utils.DisplayInfo(monNom, "Fin", "COUPURE COHÉRENTE !")
	} else {
		utils.DisplayInfo(monNom, "Fin", "Coupure non cohérente...")
	}
}
