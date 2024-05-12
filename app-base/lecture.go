package main

import (
	"fmt"
	"utils"
)

// Quand le programme n'est pas en train d'écrire, il lit
func lecture() {
	var rcvmsg string

	for {
		fmt.Scanln(&rcvmsg)

		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu")
			continue
		}
		mutex.Lock()
		if rcvmsg[0] == uint8('A') { // On traite le message s'il commence par un 'A'
			//Traitement messages sauvegarde quand la sauvegarde a été terminée
			if utils.TrouverValeur(rcvmsg[1:], "listSauvegarde") != "" {
				traiterMessageSauvegarde(rcvmsg[1:])
			} else if utils.TrouverValeur(rcvmsg[1:], "positionX") != "" {
				traiterMessagePixel(rcvmsg[1:])
			} else {
				utils.DisplayError(monNom, "lecture", "Message non supporté")
			}
		}
		rcvmsg = ""
		mutex.Unlock()
	}
}

func traiterMessagePixel(str string) {
	messagePixel := utils.StringToMessagePixel(str)
	changerPixel(messagePixel)
}

func traiterMessageSauvegarde(str string) {
	messageSauvegarde := utils.StringToMessageSauvegarde(str)
	utils.DisplayError(monNom, "lecture", "Message sauvegarde reçu : "+utils.MessageSauvegardeToString(messageSauvegarde))

	//Traitement du message de sauvegarde : enregistrement dans un fichier et notification au frontend
	//Ecrire les champs "date" et "pixels" dans un fichier
}

func changerPixel(messagePixel utils.MessagePixel) {
	//utils.DisplayError(monNom, "changerPixel", "Et là bim on change le pixel")
}
