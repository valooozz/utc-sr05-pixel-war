package main

import (
	"fmt"
	"utils"
)

// Envoi une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	mutex.Lock()
	fmt.Println(message)
	mutex.Unlock()
}

// Utile lorsque l'on doit conserver un ordre précis dans les messages (ce que ne font pas les go-routines)
func envoiSequentiel(message string) {
	fmt.Println(message)
}

// Envoie un type Message pour les applis de contrôle
func envoyerMessageControle(message utils.Message) {
	go envoyerMessage(utils.MessageToString(message))
}

// Envoie un type MessageEtat pour les applis de contrôle
func envoyerMessageEtat(messageEtat utils.MessageEtat) {
	go envoyerMessage(utils.MessageEtatToString(messageEtat))
}

// Envoie un type MessagePixel pour l'appli de base
func envoyerMessageBase(messagePixel utils.MessagePixel) {
	go envoyerMessage("A" + utils.MessagePixelToString(messagePixel))
}

// Envoie un type MessageSauvegarde pour l'appli de base
func envoyerMessageBaseSauvegarde(messageSauvegarde utils.MessageSauvegarde) {
	go envoyerMessage("A" + utils.MessageSauvegardeToString(messageSauvegarde))
}

/////////////////////
// Exclusion mutuelle
/////////////////////

// Envoie un message de SC (Requete ou Liberation) pour l'anneau
func envoyerMessageSCControle(msgSC utils.MessageExclusionMutuelle) {
	msg := utils.MessageExclusionMutuelleToString(msgSC)
	go envoyerMessage(msg)
}

// Envoie un message Accuse pour l'anneau
func envoyerMessageAccuse(msgAcc utils.MessageAccuse) {
	msg := utils.MessageAccuseToString(msgAcc)
	envoiSequentiel(msg)
}

// Envoie un message SC pour l'application de base
func envoyerMessageSCBase(msgSC utils.TypeSC) {
	msg := "A" + utils.MessageTypeSCToString(msgSC)
	envoiSequentiel(msg)
}
