package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"utils"
)

// Pour l'instant, boucle sur l'entrée standard, lit et communique le résultat à la routine d'écriture
func lecture() {
	var rcvmsg string
	var id = -1
	for {
		fmt.Scanln(&rcvmsg)
		//utils.DisplayError(monNom, "lecture", "Message reçu "+rcvmsg)

		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu -> Fin du processus")
			break
		}

		// Mise à jour de N pour un ancien
		if rcvmsg[0:3] == "CN=" {
			indiceBarre := strings.IndexByte(rcvmsg, byte('|'))
			siteDemandeur, _ := strconv.Atoi(rcvmsg[indiceBarre+1:])
			newN, _ := strconv.Atoi(rcvmsg[3:indiceBarre])
			if newN > N {
				tabSC = append(tabSC, utils.MessageExclusionMutuelle{Type: utils.Liberation, Estampille: utils.Estampille{Site: N - 1, Horloge: 0}}) //On peut se permettre de mettre l'horloge à 0 puisqu'elle va changer après l'accusé
				utils.DisplayWarning(monNom, "Arrivée", "Je mets mon N à "+strconv.Itoa(newN)+" | len(tabSC)="+strconv.Itoa(len(tabSC))+" | tabSC="+utils.TabSCToString(tabSC))
			} else {
				tabSC[siteDemandeur-1] = utils.MessageExclusionMutuelle{utils.Liberation, utils.Estampille{tabSC[siteDemandeur-1].Estampille.Site, math.MaxInt}} // ATTENTION, ZONE TRÈS SENSIBLE ICI : On ne peut pas se permettre de mettre l'horloge à 0 car elle ne va pas changer avant réception du message
				utils.DisplayWarning(monNom, "Départ", "Je mets mon N à "+strconv.Itoa(newN)+" | len(tabSC)="+strconv.Itoa(len(tabSC))+" | tabSC="+utils.TabSCToString(tabSC))
			}
			N = newN
			continue
		}

		// Mise à jour de N pour le nouveau
		if rcvmsg[0:4] == "CNN=" {
			N, _ = strconv.Atoi(rcvmsg[4:])
			utils.DisplayWarning(monNom, "lecture", "Je mets mon N à "+strconv.Itoa(N))
			tabSC = make([]utils.MessageExclusionMutuelle, N)
			for i := 0; i < len(tabSC); i++ {
				tabSC[i].Type = utils.Liberation
				tabSC[i].Estampille = utils.Estampille{Site: i, Horloge: 0}
			}
			go lancerSaveAuto()
			continue
		}

		mutex.Lock()

		// On traite uniquement les messages qui commencent par un 'C'
		if rcvmsg[0] == uint8('C') {
			rcvmsg = rcvmsg[1:]
			//utils.DisplayInfoSauvegarde("JE SUIS LÀ", "", rcvmsg)

			if utils.TrouverValeur(rcvmsg, "id") != "" { //Cas d'un message en provenance d'en bas
				id, _ = strconv.Atoi(utils.TrouverValeur(rcvmsg, "id"))
				rcvmsg = utils.TrouverValeur(rcvmsg, "message")
				//utils.DisplayWarning(monNom, "Reception", "Je vais traiter ceci :"+rcvmsg)
			}

			if rcvmsg == "sauvegarde" {
				traiterDebutSauvegarde() //OK
				// Traitement des messages de contrôle
			} else if utils.TrouverValeur(rcvmsg, "couleur") != "" {
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" {
					traiterMessagePrepost(id, rcvmsg) //OK
				} else {
					traiterMessageControle(id, rcvmsg) //OK
				}
			} else if utils.TrouverValeur(rcvmsg, "etat") != "" {
				traiterMessageEtat(id, rcvmsg) //OK
			} else if utils.TrouverValeur(rcvmsg, "siteCible") != "" {
				traiterMessageAccuse(id, rcvmsg) //OK
			} else if utils.TrouverValeur(rcvmsg, "estampilleSite") != "" {
				//utils.DisplayError(monNom, "Reception", "J'ai reçu :"+rcvmsg)
				demande := utils.StringToMessageTypeSC(rcvmsg)
				switch demande {
				case utils.Requete:
					traiterMessageRequete(id, rcvmsg) //OK
				case utils.Liberation:
					traiterMessageLiberation(id, rcvmsg) //OK
				default:
					utils.DisplayError(monNom, "lecture", "Demande de SC non supportée")
				}
			} else if utils.TrouverValeur(rcvmsg, "typeSC") != "" {
				traiterMessageSC(rcvmsg) //OK
			} else {
				traiterMessagePixel(rcvmsg) //OK
			}
		}
		mutex.Unlock()
		id = -1
	}
}

// TRAITEMENT DES CONTRÔLES NORMAUX : on extrait le pixel que l'on exploite dans l'app-base et on fait suivre l'information
// et tout cela avec les bonnes informations mises à jour dans le message : horloge, couleur
func traiterMessageControle(id int, rcvmsg string) {
	message := utils.StringToMessage(rcvmsg)

	// On traite le message uniquement s'il ne vient pas de nous
	if message.Nom == monNom {
		go envoyerMessage(toMessageIdForNet(id, ""))
		return
	}

	utils.DisplayInfo(monNom, "Controle", "Message de contrôle reçu : "+rcvmsg)

	// Extraction de la partie pixel
	messagePixel := message.Pixel

	// Première fois qu'on reçoit l'ordre de transmettre sa sauvegarde
	if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
		maCouleur = utils.Jaune

		utils.DisplayInfoSauvegarde(monNom, "Controle", "Passage en jaune")

		messageEtat := utils.MessageEtat{monEtatLocal}
		utils.DisplayInfoSauvegarde(monNom, "Controle", "Etat : "+utils.MessageEtatToString(messageEtat))
		envoyerMessageEtat(-1, messageEtat)

		// Réception d'un message prépost pas encore marqué comme prépost
	} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
		utils.DisplayInfoSauvegarde(monNom, "Controle", "Prepost détecté")
		messagePrepost := message
		messagePrepost.Prepost = true
		envoyerMessageControle(-1, messagePrepost) //?
	}

	message.Couleur = maCouleur

	// Mise à jour de l'horloge vectorielle locale et mise à jour de sa valeur dans le message également
	horlogeVectorielle = utils.MajHorlogeVectorielle(monNom, horlogeVectorielle, message.Vectorielle)
	message.Vectorielle = horlogeVectorielle

	// On met à jour l'état local
	monEtatLocal = utils.MajEtatLocal(monEtatLocal, messagePixel)
	monEtatLocal.Vectorielle = utils.CopyHorlogeVectorielle(horlogeVectorielle)

	envoyerMessageControle(id, message) // Pour la prochaine app de contrôle de l'anneau
	envoyerMessageBase(messagePixel)    // Pour l'app de base

	maCouleur = utils.Blanc //Pour multiplier les sauvegardes
}

func traiterMessagePrepost(id int, rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayInfoSauvegarde(monNom, "Prepost", "Prepost transféré : "+rcvmsg)
		go envoyerMessage(toMessageIdForNet(id, rcvmsg)) // On fait suivre le message sur l'anneau
		return
	}

	// On ajoute le message prepost à la sauvegarde générale
	message := utils.StringToMessage(rcvmsg)
	etatGlobal.ListMessagePrepost = append(etatGlobal.ListMessagePrepost, message)

	if nbEtatsAttendus == 0 {
		utils.DisplayInfoSauvegarde(monNom, "Prepost", "Fin par prepost")
		finSauvegarde()
	}
}

func traiterMessageEtat(id int, rcvmsg string) {

	if !jeSuisInitiateur {
		utils.DisplayInfoSauvegarde(monNom, "Etat", "Transfert message etat : "+rcvmsg)
		go envoyerMessage(toMessageIdForNet(id, rcvmsg))
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

	horlogeVectorielle[monNom]++

	// Mise à jour de l'état local
	monEtatLocal = utils.MajEtatLocal(monEtatLocal, messagePixel)
	monEtatLocal.Vectorielle = utils.CopyHorlogeVectorielle(horlogeVectorielle)

	message := utils.Message{messagePixel, horlogeVectorielle, monNom, maCouleur, false}
	envoyerMessageControle(-1, message)
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
	maCouleur = utils.Blanc
	jeSuisInitiateur = false
	utils.DisplayInfoSauvegarde(monNom, "Fin", "Sauvegarde complétée")

	// On affiche l'état global (liste des états locaux et liste des messages préposts)
	for _, etatLocal := range etatGlobal.ListEtatLocal {
		utils.DisplayInfoSauvegarde(monNom, "Fin", utils.EtatLocalToString(etatLocal))
	}

	for _, mp := range etatGlobal.ListEtatLocal[1].ListMessagePixel {
		utils.DisplayInfoSauvegarde(monNom, "MAJ ETAT LOCAL", "Je prends dans l'état de "+etatGlobal.ListEtatLocal[1].NomSite)
		monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, mp)
	}

	for _, mp := range etatGlobal.ListMessagePrepost {
		utils.DisplayInfoSauvegarde(monNom, "Fin", utils.MessageToString(mp))
		monEtatLocal.ListMessagePixel = append(monEtatLocal.ListMessagePixel, mp.Pixel)
	}

	// On calcule si la coupure et cohérente et on récupère sa date (max de la vectorielle de chaque site)
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

	// On met à jour l'horloge locale et le tableau de la file d'attente
	HEM++
	message := utils.MessageExclusionMutuelle{Type: demande, Estampille: utils.Estampille{Site: Site, Horloge: HEM}}
	tabSC[Site] = message

	// On transmet la Requete ou la Liberation sur l'anneau
	envoyerMessageSCControle(-1, message)
}

// APP CONTROL -> APP CONTROL
func traiterMessageRequete(id int, rcvmsg string) {
	demande := utils.StringToMessageExclusionMutuelle(rcvmsg)

	// Si le message ne vient pas de nous
	if demande.Estampille.Site != Site {

		// On met à jour l'horloge et le tableau de la file d'attente
		HEM = utils.Recaler(demande.Estampille.Horloge, HEM)
		tabSC[demande.Estampille.Site] = demande

		// On envoie un Accuse à l'émetteur de la Requete et on transmet celle-ci sur l'anneau
		envoyerMessageAccuse(-1, utils.MessageAccuse{SiteCible: demande.Estampille.Site, Estampille: utils.Estampille{Site, HEM}})
		//utils.DisplayError(monNom, "Reception", "J'envoi cet accusé :"+utils.MessageAccuseToString(utils.MessageAccuse{SiteCible: demande.Estampille.Site, Estampille: utils.Estampille{Site, HEM}}))
		envoyerMessageSCControle(id, demande)

		// On regarde si on peut accepter une SC chez nous
	} else if utils.QuestionEntreeSC(Site, tabSC) {
		utils.DisplayInfoSC(monNom, "Requete", "SC acceptée par Requete !")
		envoyerMessageSCBase(tabSC[Site].Type)
	} else {
		utils.DisplayInfoSC(monNom, "Requete", "SC refusée ! "+" La SC est occupée par C"+strconv.Itoa(utils.PlusVieilleRequeteAlive(Site, tabSC)+1)+" !")
		//utils.DisplayInfoSC(monNom, "Requete", "tabSC="+utils.TabSCToString(tabSC))
		//os.Exit(1)
	}
}

// APP CONTROL -> APP CONTROL
func traiterMessageLiberation(id int, rcvmsg string) {
	liberation := utils.StringToMessageExclusionMutuelle(rcvmsg)

	// Si le message ne vient pas de nous
	if liberation.Estampille.Site != Site {

		// On met à jour l'horloge et le tableau de la file d'attente
		HEM = utils.Recaler(liberation.Estampille.Horloge, HEM)
		tabSC[liberation.Estampille.Site] = liberation

		// On transmet le message sur l'anneau
		envoyerMessageSCControle(id, liberation)
	}

	// On regarde si on peut accepter une SC chez nous
	if utils.QuestionEntreeSC(Site, tabSC) {
		utils.DisplayInfoSC(monNom, "Liberation", "SC acceptée par Liberation !")
		envoyerMessageSCBase(tabSC[Site].Type)
	}
}

// APP CONTROL -> APP CONTROL
func traiterMessageAccuse(id int, rcvmsg string) {
	message := utils.StringToMessageAccuse(rcvmsg)

	// Si l'Accuse n'est pas pour nous, on le transmet et on quitte la fonction
	if Site != message.SiteCible {
		envoiSequentiel(toMessageIdForNet(id, rcvmsg))
		return
	}

	// Si le site qui envoie l'Accuse ne fait pas de requête, on passe son état à Accuse dans le tableau de la file d'attente
	if tabSC[message.Estampille.Site].Type != utils.Requete {
		tabSC[message.Estampille.Site] = utils.MessageExclusionMutuelle{utils.Accuse, message.Estampille}
	}
}

func lancerSaveAuto() {
	time.Sleep(1 * time.Second)
	traiterDebutSauvegarde() //Demande une sauvegarde automatiquement en arrivant
	traiterMessagePixel("/=positionX=1000/=positionY=1000/=R=0/=G=0/=B=0")
}
