package main

import (
	"fmt"
)

// Pour l'instant, boucle sur l'entrée standard, lit et communique le résultat à la routine d'écriture
func lecture() {
	var rcvmsg string

	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		if trouverValeur(rcvmsg, "horloge") != "" {
			if trouverValeur(rcvmsg, "prepost") == "true" {
				if initiateur {
					// Traiter l'ajout du message à l'état de sauvegarde
				} else {
					envoyer_message(rcvmsg)
				}
			} else {
				displayInfo("main", "Message de contrôle reçu : "+rcvmsg)
				message := StringToMessage(rcvmsg)
				K = recaler(K, message.horloge)
				messagePixel := message.pixel
				//forward à l'écriture pour l'envoi à l'app de base (avec channel définit à l'initialisation)
			}
		} else if trouverValeur(rcvmsg, "etat") != "" {
			//CAS D'UN MESSAGE DE TYPE ÉTAT DESTINÉ AU SUCCESSEUR OU PAS SI INITIATEUR => traité plus tard
			//messsageEtat := StringToMessageEtat(rcvmsg)
			//faire d'autres traitements ici
			//forward à l'écriture
		} else {
			//CAS D'UN MESSAGE DE TYPE MESSAGEPIXEL DESTINÉ AU RÉSEAU
			messagePixel := StringToMessagePixel(rcvmsg)
			stderr.Println("Message pixel reçu : ", rcvmsg)
			K++
			message := Message{messagePixel, K, maCouleur, false}
			//faire d'autres traitements ici
			//forward à l'écriture pour l'envoi sur le réseau
		}
		mutex.Unlock()
	}
}
