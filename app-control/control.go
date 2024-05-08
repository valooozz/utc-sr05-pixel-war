package main

import (
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"utils"
)

// Définition des variables
var mutex = &sync.Mutex{}
var H = 0
var horlogeVectorielle = utils.HorlogeVectorielle{}
var maCouleur = utils.Blanc
var jeSuisInitiateur = false
var monEtatLocal utils.EtatLocal
var etatGlobal utils.EtatGlobal
var monBilan = 0
var nbEtatsAttendus = 0
var nbMessagesAttendus = 0

var N = 3

var pNom = flag.String("n", "controle", "nom")
var monNom string

func main() {
	flag.Parse()
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())

	horlogeVectorielle[monNom] = 0
	monEtatLocal.NomSite = monNom
	monEtatLocal.Vectorielle = horlogeVectorielle

	/*el1 := utils.EtatLocal{monNom, utils.HorlogeVectorielle{monNom: 3}, []utils.MessagePixel{}}
	el2 := utils.EtatLocal{"C2", utils.HorlogeVectorielle{monNom: 3, "C2": 3}, []utils.MessagePixel{}}
	el3 := utils.EtatLocal{"C3", utils.HorlogeVectorielle{monNom: 3, "C2": 2, "C3": 3}, []utils.MessagePixel{}}
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, el1)
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, el2)
	etatGlobal.ListEtatLocal = append(etatGlobal.ListEtatLocal, el3)

	if utils.CoupureEstCoherente(etatGlobal) {
		fmt.Println("Coupure cohérente")
	} else {
		fmt.Println("Coupure non cohérente")
	}*/

	go lecture()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
