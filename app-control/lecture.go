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
				//TRAITEMENT DU CAS PARTICULIER DES CONTRÔLES-PRÉPOST : ceux qui sont marqués comme prépost
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" {
					if jeSuisInitiateur {
						// Traiter l'ajout du message à l'état de sauvegarde
					} else {
						go envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
					}
					//TRAITEMENT DES CONTRÔLES NORMAUX : on extrait le pixel que l'on exploite dans l'app-base et on fait suivre l'information
					// et tout cela avec les bonnes informations mises à jour dans le message : horloge, couleur
				} else {
					utils.DisplayWarning(monNom, "main", "Message de contrôle reçu : "+rcvmsg)
					//Transformation du message string en objet message => pour un meilleur traitement
					message := utils.StringToMessage(rcvmsg)

					// On traite le message uniquement s'il ne vient pas de nous
					if message.Nom != monNom {
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
								go envoyerMessage(utils.MessageToString(messagePrepost))
							}
						}

						message.Couleur = maCouleur
						go envoyerMessage(utils.MessageToString(message))                 // Pour la prochaine app de contrôle de l'anneau
						go envoyerMessage("A" + utils.MessagePixelToString(messagePixel)) // Pour l'app de base (on ajoute un 'A' au début)
					}
				}
				//TRAITEMENT DES MESSAGES D'ÉTAT
			} else if utils.TrouverValeur(rcvmsg, "etat") != "" { //
				if jeSuisInitiateur {
					// Traiter l'ajout de l'état à la sauvegarde générale
				} else {
					go envoyerMessage(rcvmsg)
				}
				//TRAITEMENT DES MESSAGES PIXELS (ENVOYÉS PAR L'APP DE BASE)
			} else {
				messagePixel := utils.StringToMessagePixel(rcvmsg)
				utils.DisplayWarning(monNom, "lecture", "Message pixel reçu : "+rcvmsg)
				H++
				//ATTENTION ICI, METTRE À JOUR L'ÉTAT GLOBAL
				message := utils.Message{messagePixel, H, monNom, maCouleur, false}
				go envoyerMessage(utils.MessageToString(message))
			}
		}
		mutex.Unlock()
	}
}
