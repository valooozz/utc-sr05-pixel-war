package main

import (
	"fmt"
	"time"
	"utils"
)

// Envoie une chaine de caractères sur la sortie standard
func envoyerMessage(message string) {
	mutex.Lock()
	msg := "C" + message
	fmt.Println(msg)
	utils.DisplayInfo(monNom, "J'envoie ", msg)
	mutex.Unlock()
}

// Utile lorsque l'on doit conserver un ordre précis dans les messages (ce que ne font pas les go-routines)
func envoiSequentiel(message string) {
	msg := "C" + message
	fmt.Println(msg)
}

// Demande à l'app de contrôle d'accéder à la section critique (Requete)
func demandeSC() {
	// On envoie un message prévenant la requête
	msg := utils.MessageTypeSCToString(utils.Requete)
	envoyerMessage(msg)

	// On bloque le programme tant qu'on n'a pas accès à la section critique
	for accesSC == false {
		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

// Prévient l'app de contrôle qu'on n'a plus besoin de la section critique (Liberation)
func relacherSC() {
	// On envoie un message pour prévenir que la section critique est libérée
	accesSC = false
	msg := utils.MessageTypeSCToString(utils.Liberation)
	envoyerMessage(msg)
}
