package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"utils"
)

// Le programme envoie périodiquement des messages sur stdout
func sendperiodic() {
	var sndmsg string
	var i int

	i = 0

	for {
		mutex.Lock()
		i = i + 1
		sndmsg = "message_" + strconv.Itoa(i) + "\n"
		fmt.Print(sndmsg)
		mutex.Unlock()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func envoyerPixel(positionX int, positionY int, rouge int, vert int, bleu int) {
	mutex.Lock()
	messagePixel := utils.MessagePixel{positionX, positionY, rouge, vert, bleu}
	fmt.Printf(utils.MessagePixelToString(messagePixel))
	mutex.Unlock()
}

// Quand le programme n'est pas en train d'écrire, il lit
func lecture() {
	var rcvmsg string

	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		utils.DisplayInfo("app-de-base", "lecture", "Réception de : "+rcvmsg)
		if rcvmsg[0] == uint8('A') { // On traite le message s'il commence par un 'A'
			messagePixel := utils.StringToMessagePixel(rcvmsg[1:])
			changerPixel(messagePixel)
		}
		mutex.Unlock()
		rcvmsg = ""
	}
}

func changerPixel(messagePixel utils.MessagePixel) {

}

var mutex = &sync.Mutex{}

func main() {

	//Création de 2 go routines qui s'exécutent en parallèle
	go sendperiodic()
	go lecture()
	//On décide de bloquer le programme principal
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
