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
		// On vérifie que le message reçu n'est pas vide
		//utils.DisplayError(monNom, "lecture", "Message reçu : "+rcvmsg)
		if rcvmsg == "" {
			utils.DisplayError(monNom, "lecture", "Message vide reçu")
			break
		}

		mutex.Lock()

		if rcvmsg[0] == uint8('B') { // On traite le message s'il commence par un 'B'
			// Traitement messages sauvegarde quand la sauvegarde a été terminée
			if utils.TrouverValeur(rcvmsg[1:], "vectorielle") != "" {
				traiterMessageSauvegarde(rcvmsg[1:])
				// Un message contenant positionX est un messagePiexel
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
	if modeDeLancement == "g" {
		wsSend(utils.MessagePixelToString(messagePixel))
	}
}

// Enregistre la sauvegarde dans un fichier
func traiterMessageSauvegarde(str string) {
	messageSauvegarde := utils.StringToMessageSauvegarde(str)
	utils.DisplayInfoSauvegarde(monNom, "lecture", "Message sauvegarde reçu : "+utils.MessageSauvegardeToString(messageSauvegarde))

	// Ecrit les champs "date" et "pixels" dans un fichier .pw
	if cheminSauvegardes[len(cheminSauvegardes)-1] != '/' {
		cheminSauvegardes = cheminSauvegardes + "/"
	}
	// On enregistre l'heure et la date de la sauvegarde dans le nom du fichier
	now := time.Now()
	fileName := now.Format("2006-01-02_15:04:05") + ".pw"
	fichier, err := os.Create(cheminSauvegardes + fileName)
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de la création du fichier :"+err.Error())
		return
	}
	writer := bufio.NewWriter(fichier)
	// Ecriture dans le fichier de l'horloge vectorielle
	_, err = writer.WriteString(utils.HorlogeVectorielleToString(messageSauvegarde.Vectorielle) + "\n")
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de l'écriture dans le fichier :"+err.Error())
		return
	}
	// Ecriture des pixels présents sur l'interface lors de la sauvegarde
	for _, mp := range messageSauvegarde.ListMessagePixel {
		if modeDeLancement == "g" {
			wsSend(utils.MessagePixelToString(mp))
		}
		_, err := writer.WriteString(utils.MessagePixelToString(mp) + "\n")
		if err != nil {
			utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors de l'écriture dans le fichier :"+err.Error())
			return
		}
	}
	if modeDeLancement == "g" {
		wsSend(lastSent)
	}
	err = writer.Flush()
	if err != nil {
		utils.DisplayError(monNom, "traiterMessageSauvegarde", "Erreur lors du vidage du tampon dans le fichier :"+err.Error())
		return
	}
	utils.DisplayInfoSauvegarde(monNom, "traiterMessageSauvegarde", "Écriture dans le fichier terminée avec succès.")
	fichier.Close()
}

// On a reçu une validation de notre app de contrôle pour accéder à la section critique
func traiterMessageTypeSC() {
	accesSC = true
}
