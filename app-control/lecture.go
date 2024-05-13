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
			} else if utils.TrouverValeur(rcvmsg, "siteCible") != "" {
				traiterMessageAccuse(rcvmsg)
			} else if utils.TrouverValeur(rcvmsg, "estampilleSite") != "" {
				demande := utils.StringToMessageTypeSC(rcvmsg)
				switch demande {
				case utils.Requete:
					traiterMessageRequete(rcvmsg)
				case utils.Liberation:
					traiterMessageLiberation(rcvmsg)
				default:
					utils.DisplayError(monNom, "lecture", "Demande de SC non supportée")
				}
			} else if utils.TrouverValeur(rcvmsg, "typeSC") != "" {
				traiterMessageSC(rcvmsg)
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

		utils.DisplayInfoSauvegarde(monNom, "Controle", "Passage en jaune")

		messageEtat := utils.MessageEtat{monEtatLocal}
		utils.DisplayInfoSauvegarde(monNom, "Controle", "Etat : "+utils.MessageEtatToString(messageEtat))
		envoyerMessageEtat(messageEtat)

		// Réception d'un message prépost pas encore marqué comme prépost
	} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
		utils.DisplayInfoSauvegarde(monNom, "Controle", "Prepost détecté")
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
		utils.DisplayInfoSauvegarde(monNom, "Prepost", "Prepost transféré : "+rcvmsg)
		go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
		return
	}

	// On ajoute le message prepost à la sauvegarde générale
	message := utils.StringToMessage(rcvmsg)
	etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)

	// Ça normalement on y arrive jamais, mais je le laisse au cas où ?
	if nbEtatsAttendus == 0 {
		utils.DisplayInfoSauvegarde(monNom, "Prepost", "Fin par prepost")
		finSauvegarde()
	}
}

func traiterMessageEtat(rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayInfoSauvegarde(monNom, "Etat", "Transfert message etat : "+rcvmsg)
		go envoyerMessage(rcvmsg)
		return
	}

	messageEtat := utils.StringToMessageEtat(rcvmsg)
	utils.DisplayInfoSauvegarde(monNom, "Etat", "MessageEtat recu")

	// On ajoute l'état local reçu à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, utils.CopyEtatLocal(messageEtat.EtatLocal))

	nbEtatsAttendus--

	utils.DisplayInfoSauvegarde(monNom, "Etat", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus))
	if nbEtatsAttendus == 0 {
		utils.DisplayInfoSauvegarde(monNom, "Etat", "Fin par etat")
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

	utils.DisplayInfoSauvegarde(monNom, "DebutSauv", "nbEtatsAttendus="+strconv.Itoa(nbEtatsAttendus))

	// On ajoute l'état local à la sauvegarde générale
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, utils.CopyEtatLocal(monEtatLocal))
}

func finSauvegarde() {
	utils.DisplayInfoSauvegarde(monNom, "Fin", "Sauvegarde complétée")
	for _, etatLocal := range etatGlobal.ListEtatLocal {
		utils.DisplayInfoSauvegarde(monNom, "Fin", utils.EtatLocalToString(etatLocal))
	}
	for _, mp := range etatGlobal.ListMessagePrepost {
		utils.DisplayInfoSauvegarde(monNom, "Fin", utils.MessageToString(mp))
	}

	coherente, maxVectorielle := utils.CoupureEstCoherente(etatGlobal)

	if coherente {
		utils.DisplayInfoSauvegarde(monNom, "Fin", "COUPURE COHÉRENTE !")
		configurationGlobale := utils.ReconstituerCarte(etatGlobal)
		messageSauvegarde := utils.MessageSauvegarde{ListMessagePixel: configurationGlobale, Vectorielle: maxVectorielle}
		envoyerMessageBaseSauvegarde(messageSauvegarde)
	} else {
		utils.DisplayError(monNom, "Fin", "Coupure non cohérente...")
	}
}

/////////////////////
// Exclusion mutuelle
/////////////////////

// APP BASE -> APP CONTROLE
func traiterMessageSC(rcvmsg string) {
	demande := utils.StringToMessageTypeSC(rcvmsg)
	var typeToString string
	if demande == utils.Requete {
		typeToString = "REQUÊTE d'accès à la section critique"
	} else {
		typeToString = "LIBÉRATION de l'accès à la section critique"
	}
	utils.DisplayInfoSC(monNom, "SC", "A"+strconv.Itoa(Site+1)+" envoie : "+typeToString)
	HEM++
	message := utils.MessageExclusionMutuelle{Type: demande, Estampille: utils.Estampille{Site: Site, Horloge: HEM}}
	tabSC[Site] = message
	envoyerMessageSCControle(message)
}

// APP CONTROL -> APP CONTROL
func traiterMessageRequete(rcvmsg string) {
	demande := utils.StringToMessageExclusionMutuelle(rcvmsg)

	if demande.Estampille.Site != Site {
		//utils.DisplayInfoSC(monNom, "Requete", "Message de requete reçu et forwardé : "+rcvmsg)
		HEM = utils.Recaler(demande.Estampille.Horloge, HEM)
		tabSC[demande.Estampille.Site] = demande
		envoyerMessageAccuse(utils.MessageAccuse{SiteCible: demande.Estampille.Site, Estampille: utils.Estampille{Site, HEM}})
		envoyerMessageSCControle(demande)
	}

	if utils.QuestionEntreeSC(Site, tabSC) {
		utils.DisplayInfoSC(monNom, "Requete", "SC acceptée !")
		envoyerMessageSCBase(tabSC[Site].Type)
	}
}

// APP CONTROL -> APP CONTROL
func traiterMessageLiberation(rcvmsg string) {
	liberation := utils.StringToMessageExclusionMutuelle(rcvmsg)

	if liberation.Estampille.Site != Site {
		//utils.DisplayInfoSC(monNom, "Liberation", "Message de liberation reçu et forwardé : "+rcvmsg)
		HEM = utils.Recaler(liberation.Estampille.Horloge, HEM)
		tabSC[liberation.Estampille.Site] = liberation
		envoyerMessageSCControle(liberation)
	}

	if utils.QuestionEntreeSC(Site, tabSC) {
		utils.DisplayInfoSC(monNom, "Liberation", "SC acceptée !")
		envoyerMessageSCBase(tabSC[Site].Type)
	}
}

func traiterMessageAccuse(rcvmsg string) {
	message := utils.StringToMessageAccuse(rcvmsg)

	if Site != message.SiteCible {
		//utils.DisplayInfoSC(monNom, "Accuse", "Message d'accusé reçu et forwardé : "+rcvmsg)
		envoiSequentiel(rcvmsg)
		return
	}

	//utils.DisplayInfoSC(monNom, "Accuse", "Message d'accusé pour moi")
	tabSC[message.Estampille.Site] = utils.MessageExclusionMutuelle{utils.Accuse, message.Estampille}
}
