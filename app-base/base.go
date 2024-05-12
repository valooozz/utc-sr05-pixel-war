package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func fiveSecondsSnapshot(n int) {
	time.Sleep(time.Duration(n) * time.Second)
	mutex.Lock()
	fmt.Println("sauvegarde")
	mutex.Unlock()
}

// Le programme envoie périodiquement des messages sur stdout
func sendperiodic() {
	val, _ := strconv.Atoi(monNom[1:2])
	for i := 0; i < 8; i++ {
		demandeSC()
		//updateMatriceFront()
		envoyerPixel(i, i, 255, val, 0)
		relacherSC()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

var mutex = &sync.Mutex{}
var pNom = flag.String("n", "base", "nom")
var pPath = flag.String("p", "./sauvegardes", "path")
var monNom string
var cheminSauvegardes string
var accesSC = false

func main() {
	flag.Parse()
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())
	cheminSauvegardes = *pPath

	//if monNom[0:2] == "A1" {
	//	go fiveSecondsSnapshot(5)
	//}

	//Création de 2 go routines qui s'exécutent en parallèle
	//|| monNom[0:2] == "A2"
	//|| monNom[0:2] == "A2" || monNom[0:2] == "A3"
	if monNom[0:2] == "A1" || monNom[0:2] == "A2" || monNom[0:2] == "A3" {
		go sendperiodic()
	}
	go lecture()
	//On décide de bloquer le programme principal
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
