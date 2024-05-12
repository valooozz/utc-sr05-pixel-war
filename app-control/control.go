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
var nbEtatsAttendus = 0

var N = 3
var tabSC = make([]utils.MessageExclusionMutuelle, N)

var pNom = flag.String("n", "controle", "nom")
var monNom string
var Site int

var HEM = 0

func main() {
	flag.Parse()
	Site = utils.InitialisationNumSite(*pNom) - 1
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())

	for i := 0; i < len(tabSC); i++ {
		tabSC[i].Type = utils.Liberation
		tabSC[i].Estampille = utils.Estampille{Site: i, Horloge: 0}
	}

	horlogeVectorielle[monNom] = 0
	monEtatLocal.NomSite = monNom
	monEtatLocal.Vectorielle = horlogeVectorielle

	go lecture()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
