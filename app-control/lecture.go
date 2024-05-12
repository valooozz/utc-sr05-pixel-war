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
			utils.DisplayError(monNom, "lecture", "Message vide reçu -> Fin du processus")
			break
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

	utils.DisplayInfo(monNom, "Controle", "Message de contrôle reçu : "+rcvmsg)

	// Extraction de la partie pixel
	messagePixel := message.Pixel

	// Recalage de l'horloge locale et mise à jour de sa valeur dans le message également
	H = utils.Recaler(H, message.Horloge)
	message.Horloge = H

	// Première fois qu'on reçoit l'ordre de transmettre sa sauvegarde
	if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
		maCouleur = utils.Jaune

		utils.DisplayWarning(monNom, "Controle", "Passage en jaune")

		messageEtat := utils.MessageEtat{monEtatLocal}
		utils.DisplayInfo(monNom, "Controle", "Etat : "+utils.MessageEtatToString(messageEtat))
		envoyerMessageEtat(messageEtat)

		// Réception d'un message prépost pas encore marqué comme prépost
	} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
		utils.DisplayWarning(monNom, "Controle", "Prepost")
		messagePrepost := message
		messagePrepost.Prepost = true
		envoyerMessageControle(messagePrepost)
	}

	message.Couleur = maCouleur

	// Mise à jour de l'horloge vectorielle locale et mise à jour de sa valeur dans le message également
	horlogeVectorielle = utils.MajHorlogeVectorielle(monNom, horlogeVectorielle, message.Vectorielle)
	message.Vectorielle = horlogeVectorielle

	// On met à jour l'état local
	monEtatLocal = utils.MajEtatLocal(monEtatLocal, messagePixel)
	monEtatLocal.Vectorielle = utils.CopyHorlogeVectorielle(horlogeVectorielle)

	envoyerMessageControle(message)  // Pour la prochaine app de contrôle de l'anneau
	envoyerMessageBase(messagePixel) // Pour l'app de base
}

func traiterMessagePrepost(rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayInfo(monNom, "Prepost", "Prepost transféré : "+rcvmsg)
		go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
		return
	}

	// On ajoute le message prepost à la sauvegarde générale
	message := utils.StringToMessage(rcvmsg)
	etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)

	// Ça normalement on y arrive jamais, mais je le laisse au cas où ?
	if nbEtatsAttendus == 0 {
		utils.DisplayInfo(monNom, "Prepost", "Fin par prepost")
		finSauvegarde()
	}
}

func traiterMessageEtat(rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayInfo(monNom, "Etat", "Transfert message etat : "+rcvmsg)
		go envoyerMessage(rcvmsg)
		return
	}

	messageEtat := utils.StringToMessageEtat(rcvmsg)
	utils.DisplayInfo(monNom, "Etat", "MessageEtat recu")

	// On ajoute l'état local reçu à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, utils.CopyEtatLocal(messageEtat.EtatLocal))

	nbEtatsAttendus--

	utils.DisplayWarning(monNom, "Etat", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus))
	if nbEtatsAttendus == 0 {
		utils.DisplayWarning(monNom, "Etat", "Fin par etat")
		finSauvegarde()
	}
}

func traiterMessagePixel(rcvmsg string) {
	utils.DisplayInfo(monNom, "Pixel", "Message pixel reçu : "+rcvmsg)

	messagePixel := utils.StringToMessagePixel(rcvmsg)
	H++

	horlogeVectorielle[monNom]++

	// Mise à jour de l'état local
	monEtatLocal = utils.MajEtatLocal(monEtatLocal, messagePixel)
	monEtatLocal.Vectorielle = utils.CopyHorlogeVectorielle(horlogeVectorielle)

	message := utils.Message{messagePixel, H, horlogeVectorielle, monNom, maCouleur, false}
	envoyerMessageControle(message)
}

func traiterDebutSauvegarde() {
	maCouleur = utils.Jaune
	jeSuisInitiateur = true
	nbEtatsAttendus = N - 1

	utils.DisplayWarning(monNom, "DebutSauv", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus))

	// On ajoute l'état local à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, utils.CopyEtatLocal(monEtatLocal))
}

func finSauvegarde() {
	utils.DisplayWarning(monNom, "Fin", "Sauvegarde complétée")
	for _, etatLocal := range etatGlobal.ListEtatLocal {
		utils.DisplayWarning(monNom, "Fin", utils.EtatLocalToString(etatLocal))
	}
	for _, mp := range etatGlobal.ListMessagePrepost {
		utils.DisplayWarning(monNom, "Fin", utils.MessageToString(mp))
	}

	coherente, maxVectorielle := utils.CoupureEstCoherente(etatGlobal)

	if coherente {
		utils.DisplayWarning(monNom, "Fin", "COUPURE COHÉRENTE !")
		configurationGlobale := utils.ReconstituerCarte(etatGlobal)
		messageSauvegarde := utils.MessageSauvegarde{ListMessagePixel: configurationGlobale, Vectorielle: maxVectorielle}
		envoyerMessageBaseSauvegarde(messageSauvegarde)
	} else {
		utils.DisplayError(monNom, "Fin", "Coupure non cohérente...")
	}
}
