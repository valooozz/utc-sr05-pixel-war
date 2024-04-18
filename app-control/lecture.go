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
		if rcvmsg[0] != uint8('A') { // On traite uniquement les messages qui ne commencent pas par un 'A'
			if utils.TrouverValeur(rcvmsg, "horloge") != "" { // Message de contrôle d'ajout de pixel
				if utils.TrouverValeur(rcvmsg, "prepost") == "true" { // Si c'est un message prepost
					if jeSuisInitiateur {
						// Traiter l'ajout du message à l'état de sauvegarde
					} else {
						envoyerMessage(rcvmsg) // On fait suivre le message sur l'anneau
					}
				} else {
					utils.DisplayInfo(monNom, "main", "Message de contrôle reçu : "+rcvmsg)
					message := utils.StringToMessage(rcvmsg)

					if message.Nom != monNom { // On traite le message uniquement s'il ne vient pas de nous
						messagePixel := message.Pixel
						K = utils.Recaler(K, message.Horloge)
						message.Horloge = K

						if message.Couleur == utils.Jaune && maCouleur == utils.Blanc {
							maCouleur = utils.Jaune
							messageEtat := utils.MessageEtat{list.List(monEtatLocal), monBilan}
							envoyerMessage(utils.MessageEtatToString(messageEtat))
						} else if message.Couleur == utils.Blanc && maCouleur == utils.Jaune {
							if jeSuisInitiateur {
								// Ajouter le message reçu à la sauvegarde générale
							} else {
								messagePrepost := message
								messagePrepost.Prepost = true
								envoyerMessage(utils.MessageToString(messagePrepost))
							}
						}

						message.Couleur = maCouleur
						envoyerMessage(utils.MessageToString(message))                 // Pour la prochaine app de contrôle de l'anneau
						envoyerMessage("A" + utils.MessagePixelToString(messagePixel)) // Pour l'app de base (on ajoute un 'A' au début)
					}
				}
			} else if utils.TrouverValeur(rcvmsg, "etat") != "" {
				if jeSuisInitiateur {
					// Traiter l'ajout de l'état à la sauvegarde générale
				} else {
					envoyerMessage(rcvmsg)
				}
			} else { // MessagePixel envoyé par l'app de base
				messagePixel := utils.StringToMessagePixel(rcvmsg)
				utils.DisplayInfo(monNom, "lecture", "Message pixel reçu : "+rcvmsg)
				K++
				message := utils.Message{messagePixel, K, monNom, maCouleur, false}
				envoyerMessage(utils.MessageToString(message))
			}
		}
		mutex.Unlock()
	}
}
