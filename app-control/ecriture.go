package main

import (
	"fmt"
	"utils"
)

// Envoi une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	mutex.Lock()
	//utils.DisplayError(monNom, "envoyer", message)
	fmt.Println(message)
	mutex.Unlock()
}

// Utile lorsque l'on doit conserver un ordre précis dans les messages (ce que ne font pas les go-routines)
func envoiSequentiel(message string) {
	//utils.DisplayError(monNom, "envoyer", message)
	fmt.Println(message)
}

// Envoie un type Message pour les applis de contrôle
func envoyerMessageControle(id int, message utils.Message) {
	go envoyerMessage(toMessageIdForNet(id, utils.MessageToString(message)))
}

// Envoie un type MessageEtat pour les applis de contrôle
func envoyerMessageEtat(id int, messageEtat utils.MessageEtat) {
	go envoyerMessage(toMessageIdForNet(id, utils.MessageEtatToString(messageEtat)))
}

// Envoie un type MessagePixel pour l'appli de base
func envoyerMessageBase(messagePixel utils.MessagePixel) {
	go envoyerMessage("B" + utils.MessagePixelToString(messagePixel))
}

// Envoie un type MessageSauvegarde pour l'appli de base
func envoyerMessageBaseSauvegarde(messageSauvegarde utils.MessageSauvegarde) {
	go envoyerMessage("B" + utils.MessageSauvegardeToString(messageSauvegarde))
}

/////////////////////
// Exclusion mutuelle
/////////////////////

// Envoie un message de SC (Requete ou Liberation) pour l'anneau
func envoyerMessageSCControle(id int, msgSC utils.MessageExclusionMutuelle) {
	msg := utils.MessageExclusionMutuelleToString(msgSC)
	go envoyerMessage(toMessageIdForNet(id, msg))
}

// Envoie un message Accuse pour l'anneau
func envoyerMessageAccuse(id int, msgAcc utils.MessageAccuse) {
	msg := utils.MessageAccuseToString(msgAcc)
	envoiSequentiel(toMessageIdForNet(id, msg))
}

// Envoie un message SC pour l'application de base
func envoyerMessageSCBase(msgSC utils.TypeSC) {
	msg := "B" + utils.MessageTypeSCToString(msgSC)
	envoiSequentiel(msg)
}

/////////////////////
// Communication avec l'app NET
/////////////////////

// Conversion en message destiné à l'app NET
func toMessageIdForNet(id int, msg string) string {
	messageId := utils.MessageId{Id: id, Message: msg}
	return "N" + utils.MessageIdToString(messageId)
}
