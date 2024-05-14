package main

import (
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"utils"
)

func nSecondsSnapshot(n int) {
	time.Sleep(time.Duration(n) * time.Second)
	mutex.Lock()
	envoiSequentiel("sauvegarde")
	mutex.Unlock()
}

// Le programme envoie périodiquement des messages sur stdout
func sendPeriodic(nbMessages int, slower bool) {
	val, _ := strconv.Atoi(monNom[1:2])
	for i := 0; i < nbMessages; i++ {
		demandeSC()
		//Le slower permet créer une différence de vitesse entre les sites et accentue la dispute pour la section critique
		//Ici que pour les 2 premiers sites
		if slower {
			if monNom[0:2] == "A1" {
				time.Sleep(time.Duration(3) * time.Second)
			}
			if monNom[0:2] == "A2" {
				time.Sleep(time.Duration(1) * time.Second)
			}
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

var pMode = flag.String("m", "a", "mode") //"t" ou "g" ou "a" pour terminal ou graphique ou "automatique
var pAutoSave = flag.Int("s", -1, "sauvegarde automatique au bout de n secondes (>=0)")

func main() {
	flag.Parse()
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())
	cheminSauvegardes = *pPath
	modeDeLancement := *pMode
	autoSave := *pAutoSave

	//Si l'option m == "g" on lance l'interface graphique, sinon le mode terminal ou automatique
	if modeDeLancement == "g" {
		//LANCEMENT DE L'INTERFACE GRAPHIQUE DANS UNE GO ROUTINE : car elle vient remplacer sendPeriodic
		//C'est à l'interface d'utiliser demandeSC(), envoyerPixel() et relacherSC() pour faire ses transactions
		//Ces fonctions se chargent de l'aspect mutex dans l'app de base, de l'aspect section critique également grace au booléen dédié
		//[Lancement ici]
	} else if modeDeLancement == "t" {
		//LANCEMENT DU MODE TERMINAL
		//On lance une sauvegarde au bout de n secondes
		if autoSave >= 0 {
			go nSecondsSnapshot(autoSave)
		}

		//On lance une fonction d'envoi périodique sur la diagonale (20 messages)
		go sendPeriodic(20, false)
	} else if modeDeLancement == "a" {
		//LANCEMENT DU MODE AUTOMATIQUE
		//On lance le snapshot sur A1 au bout de 7 secondes (A1 doit être en mode automatique biensûr)
		if monNom[0:2] == "A1" {
			go nSecondsSnapshot(7)
		}

		//On lance un envoi automatique périodique sur la diagonale sur les 2 premiers/seuls sites (ils doivent exister sous ce nom biensûr)
		if monNom[0:2] == "A1" || monNom[0:2] == "A2" {
			go sendPeriodic(20, true)
		}
	}
	go lecture()
	//On décide de bloquer le programme principal
	for {
		time.Sleep(time.Duration(60) * time.Second)
	}
}
