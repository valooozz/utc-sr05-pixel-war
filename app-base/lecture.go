package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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
			} else if utils.TrouverValeur(rcvmsg[1:], "typeSC") != "" {
				traiterMessageTypeSC()
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
	utils.DisplayInfoSauvegarde(monNom, "lecture", "Message sauvegarde reçu : "+utils.MessageSauvegardeToString(messageSauvegarde))

	//Traitement du message de sauvegarde : enregistrement dans un fichier et notification au frontend

	//Ecrire les champs "date" et "pixels" dans un fichier .pw
	if cheminSauvegardes[len(cheminSauvegardes)-1] != '/' {
		cheminSauvegardes = cheminSauvegardes + "/"
	}
	now := time.Now()
	fileName := now.Format("2006-01-02_15:04:05") + ".pw"
	fichier, err := os.Create(cheminSauvegardes + fileName)
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de la création du fichier :"+err.Error())
		return
	}
	writer := bufio.NewWriter(fichier)

	_, err = writer.WriteString(utils.HorlogeVectorielleToString(messageSauvegarde.Vectorielle) + "\n")
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de l'écriture dans le fichier :"+err.Error())
		return
	}

	for _, mp := range messageSauvegarde.ListMessagePixel {
		_, err := writer.WriteString(utils.MessagePixelToString(mp) + "\n")
		if err != nil {
			utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de l'écriture dans le fichier :"+err.Error())
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors du vidage du tampon dans le fichier :"+err.Error())
		return
	}
	utils.DisplayInfoSauvegarde(monNom, "traiterMessageSauvegarde", "Écriture dans le fichier terminée avec succès.")
	fichier.Close()

	//Notification au frontend
	//A FAIRE UNE FOIS LE FRONTEND TERMINÉ
}

func traiterMessageTypeSC() {
	//mettre le boolean acces à true
	accesSC = true
}

func changerPixel(messagePixel utils.MessagePixel) {
	//utils.DisplayError(monNom, "changerPixel", "Et là bim on change le pixel")
}
