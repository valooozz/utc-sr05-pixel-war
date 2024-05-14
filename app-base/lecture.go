package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
	"utils"
)

// Quand le programme n'est pas en train d'écrire, il lit
func lecture(game *Game) {
	utils.DisplayError(monNom, "lecture", "Lecture ici lancée")
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
				traiterMessagePixel(rcvmsg[1:], game)
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

func traiterMessagePixel(str string, game *Game) {
	if game.Matrix != nil {
		messagePixel := utils.StringToMessagePixel(str)
		changerPixel(messagePixel, game)
	} else {
		//utils.DisplayInfo(monNom, "MsgPix", "Je mets à jour avec : "+str)
	}
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

func changerPixel(messagePixel utils.MessagePixel, game *Game) {
	messageString := utils.MessagePixelToString(messagePixel)
	cr, _ := strconv.Atoi(utils.TrouverValeur(messageString, "R"))
	cb, _ := strconv.Atoi(utils.TrouverValeur(messageString, "B"))
	cg, _ := strconv.Atoi(utils.TrouverValeur(messageString, "G"))
	x, _ := strconv.Atoi(utils.TrouverValeur(messageString, "positionX"))
	y, _ := strconv.Atoi(utils.TrouverValeur(messageString, "positionY"))

	game.UpdateMatrix(x, y, uint8(cr), uint8(cg), uint8(cb))
}
