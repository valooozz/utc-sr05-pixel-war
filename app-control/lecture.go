package main

import (
	"fmt"
	"log"
	"os"
)

func recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}

// Pour l'instant, boucle sur l'entrée standard, lit et communique le résultat à la routine d'écriture
func lecture() {
	var rcvmsg string
	l := log.New(os.Stderr, "", 0)

	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		if trouverValeur(rcvmsg, "horloge") != "" {
			//CAS D'UN MESSAGE DE TYPE MESSAGE DESTINÉ A L'APP DE BASE OU PAS SI PREPOST
			message := StringToMessage(rcvmsg)
			K = recaler(K, message.horloge)
			messagePixel := message.pixel
			l.Println("Message de contrôle reçu : ", rcvmsg)
			//faire d'autres traitements ici
			//forward à l'écriture pour l'envoi à l'app de base (avec channel définit à l'initialisation)
		} else if trouverValeur(rcvmsg, "etat") != "" {
			//CAS D'UN MESSAGE DE TYPE ÉTAT DESTINÉ AU SUCCESSEUR OU PAS SI INITIATEUR => traité plus tard
			//messsageEtat := StringToMessageEtat(rcvmsg)
			//faire d'autres traitements ici
			//forward à l'écriture
		} else {
			//CAS D'UN MESSAGE DE TYPE MESSAGEPIXEL DESTINÉ AU RÉSEAU
			messagePixel := StringToMessagePixel(rcvmsg)
			l.Println("Message pixel reçu : ", rcvmsg)
			K++
			message := Message{messagePixel, K, maCouleur, false}
			//faire d'autres traitements ici
			//forward à l'écriture pour l'envoi sur le réseau
		}
		mutex.Unlock()
	}
}
