package main

import (
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"utils"
)

func envoyerPixel(positionX int, positionY int, rouge int, vert int, bleu int) {
	// On envoie un message contenant le pixel posé à l'app de contrôle
	messagePixel := utils.MessagePixel{positionX, positionY, rouge, vert, bleu}
	envoyerMessage(utils.MessagePixelToString(messagePixel))
}

func nSecondsSnapshot(n int) {
	// On envoie un message de sauvegarde automatiquement après n secondes
	time.Sleep(time.Duration(n) * time.Second)
	mutex.Lock()
	envoiSequentiel("sauvegarde")
	mutex.Unlock()
}

// Le programme envoie périodiquement des messages sur stdout
func sendPeriodic(nbMessages int, slower bool) {
	val, _ := strconv.Atoi(monNom[1:2])
	time.Sleep(time.Duration(30) * time.Second)
	for i := 0; i < nbMessages; i++ {
		// On demande l'accès à la section critique
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
			if monNom[0:2] == "A3" {
				time.Sleep(time.Duration(2) * time.Second)
			}
			if monNom[0:2] == "A4" {
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}
		// Variable val permet d'identifier le site initiateur du message
		envoyerPixel(i, i, 255, val, 0)
		// On libère la section critique
		relacherSC()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
	utils.DisplayWarning(monNom, "sendPeriodic", "SEND PERIODIC FINIT")
}

// Variables globales de répartition
var mutex = &sync.Mutex{}
var pNom = flag.String("n", "base", "nom")
var pPath = flag.String("p", "./sauvegardes", "path")
var monNom string
var cheminSauvegardes string
var accesSC = false // false si on n'est pas en section critique, true si est en section critique

// Variables globales d'utilisation
var pMode = flag.String("m", "a", "mode") //"g" ou "a" pour graphique ou automatique
var pPort = flag.Int("port", 4444, "n° de port")
var pAddr = flag.String("addr", "localhost", "nom/adresse machine")
var modeDeLancement string

func main() {
	flag.Parse()
	// On initialise le nom de l'application et le chemin vers le fichier de sauvegarde
	// ainsi que le mode de lancement
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())
	cheminSauvegardes = *pPath
	modeDeLancement = *pMode
	port := *pPort
	addr := *pAddr

	//Si l'option m == "g" on lance l'interface graphique, sinon le mode terminal ou automatique
	if modeDeLancement == "g" {
		lancementModeGraphique(strconv.Itoa(port), addr)
	} else {
		lancementModeAutomatique()
	}
}

func lancementModeGraphique(port string, addr string) {
	// On lance une goroutine de lecture des messages ainsi que l'interface graphique
	go lecture()
	launchServer(port, addr)
	//ici potentiellement lancer un client automatiquement ou dans le script
}

func lancementModeAutomatique() {
	//On lance le snapshot sur A1 au bout de 10 secondes (A1 doit être en mode automatique biensûr)
	//if monNom[0:2] == "A1" {
	//	go nSecondsSnapshot(10)
	//}

	//On lance un envoi automatique périodique sur la diagonale sur les 2 premiers/seuls sites (ils doivent exister sous ce nom biensûr)
	if monNom[0:2] == "A1" || monNom[0:2] == "A2" || monNom[0:2] == "A3" || monNom[0:2] == "A4" {
		go sendPeriodic(20, true)
	}

	go lecture()
	//On décide de bloquer le programme principal
	for {
		time.Sleep(time.Duration(60) * time.Second)
	}
}
