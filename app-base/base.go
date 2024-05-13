package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
	"utils"
)

func nSecondsSnapshot(n int) {
	time.Sleep(time.Duration(n) * time.Second)
	mutex.Lock()
	fmt.Println("sauvegarde")
	mutex.Unlock()
}

// Le programme envoie périodiquement des messages sur stdout
func sendPeriodic() {
	val, _ := strconv.Atoi(monNom[1:2])
	for i := 0; i < 20; i++ {
		demandeSC()
		//updateMatriceFront()
		if monNom[0:2] == "A1" {
			time.Sleep(time.Duration(3) * time.Second)
		}
		if monNom[0:2] == "A2" {
			time.Sleep(time.Duration(2) * time.Second)
		}
		envoyerPixel(i, i, 255, val, 0)
		relacherSC()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
	utils.DisplayWarning(monNom, "sendPeriodic", "SEND PERIODIC FINIT")
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

	//TEST DE SAUVEGARDE : on envoie une sauvegarde au bout de n secondes
	//if monNom[0:2] == "A1" {
	//	go nSecondsSnapshot(20)
	//}

	//TEST EXCLUSION MUTUELLE : lancement de plusieurs sites dont 1 plus lent quand il a l'accès à la SC
	//Création de 2 go routines qui s'exécutent en parallèle
	//|| monNom[0:2] == "A2"
	//|| monNom[0:2] == "A2" || monNom[0:2] == "A3"
	if monNom[0:2] == "A1" || monNom[0:2] == "A2" || monNom[0:2] == "A3" {
		go sendPeriodic()
	}
	go lecture()
	//On décide de bloquer le programme principal
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
