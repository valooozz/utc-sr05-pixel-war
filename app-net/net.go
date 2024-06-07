package main

import (
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"utils"
)

var mutex = &sync.Mutex{}
var pNom = flag.String("n", "controle", "nom")
var pRoutage = flag.String("r", "", "routage")
var pNbsites = flag.Int("nbsites", 3, "nom")
var pPort = flag.Int("port", 4444, "n° de port")
var pAddr = flag.String("addr", "localhost", "nom/adresse machine")

var N int         // Nombre de sites dans le réseau
var monNom string // Nom du site (option -n + pid)
var monNum int
var headers = make(map[string]utils.Header)
var siteIdCpt = 0
var tableDeRoutage = make(utils.TableDeRoutage, 0)

var monEtat string
var monParent = 0
var nbVoisinsAttendus int
var monElu = N * 100
var demande utils.Demande

var pVoisins = flag.Int("v", 0, "nombre de voisins")
var pEtatDepart = flag.String("e", "actif", "état de départ")
var pCible = flag.Int("c", 0, "site cible de la demande de raccord")
var pTimer = flag.Int("t", 0, "timer avant de rejoindre le réseau")
var pQuit = flag.Int("q", -1, "timer avant de quitter le réseau")

func main() {
	flag.Parse()
	N = *pNbsites
	monNomBrut := *pNom
	monNom = monNomBrut + "-" + strconv.Itoa(os.Getpid())
	monNum, _ = strconv.Atoi(monNomBrut[1:])
	nbVoisinsAttendus = *pVoisins
	monEtat = *pEtatDepart
	demande = utils.Demande{0, 0}
	chaineRoutage := *pRoutage
	tdr := utils.StringToTableDeRoutage(chaineRoutage)
	for _, route := range tdr {
		tableDeRoutage = append(tableDeRoutage, route)
	}

	port := *pPort
	addr := *pAddr
	// On lance une go-routine pour écouter les messages entrants sur l'entrée standard
	go lecture()
	launchServer(strconv.Itoa(port), addr)
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
