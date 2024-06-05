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
var monNom string // Nom du site (option -n + pid)
var headers = make(map[string]utils.Header)
var siteIdCpt = 0

func main() {
	flag.Parse()
	monNom = *pNom + "-" + strconv.Itoa(os.Getpid())

	// On lance une go-routine pour écouter les messages entrants sur l'entrée standard
	go lecture()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
