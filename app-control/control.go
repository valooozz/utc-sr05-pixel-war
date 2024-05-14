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
var pNom = flag.String("n", "controle", "nom")
var monNom string // Nom du site (option -n + pid)
var Site int      // Numéro du site
var N = 3         // Nombre de sites dans le réseau

var horlogeVectorielle = utils.HorlogeVectorielle{}
var maCouleur = utils.Blanc
var jeSuisInitiateur = false
var monEtatLocal utils.EtatLocal
var etatGlobal utils.EtatGlobal
var nbEtatsAttendus = 0

var HEM = 0                                           // Horloge Exclusion Mutuelle
var tabSC = make([]utils.MessageExclusionMutuelle, N) // Tableau utilisé par la file d'attente répartie afin de gérer les sections critiques

func main() {

	// On initialise le nom et le numéro du site
	flag.Parse()
	Site = utils.InitialisationNumSite(*pNom) - 1
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())

	// On initialise le tableau de la file d'attente répartie avec Liberation partout
	for i := 0; i < len(tabSC); i++ {
		tabSC[i].Type = utils.Liberation
		tabSC[i].Estampille = utils.Estampille{Site: i, Horloge: 0}
	}

	// On initialise l'horloge vectorielle avec le site local pour l'instant
	horlogeVectorielle[monNom] = 0

	// On initialise l'état local
	monEtatLocal.NomSite = monNom
	monEtatLocal.Vectorielle = horlogeVectorielle

	// On lance une go-routine pour écouter les messages entrants sur l'entrée standard
	go lecture()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
