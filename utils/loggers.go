package utils

import (
	"log"
	"os"
)

//Définition des couleurs/constantes

var rouge string = "\033[1;31m"
var orange string = "\033[1;33m"
var vert string = "\033[1;32m"
var raz string = "\033[0;00m"

//Définition des fonctions de différentes displays

var stderr = log.New(os.Stderr, "", 0)

func DisplayInfo(monNom string, where string, what string) {
	stderr.Printf("%s%s + %-8.8s : %s\n%s", vert, monNom, where, what, raz)
}

func DisplayWarning(monNom string, where string, what string) {
	stderr.Printf("%s%s * %-8.8s : %s\n%s", orange, monNom, where, what, raz)
}

func DisplayError(monNom string, where string, what string) {
	stderr.Printf("%s%s ! %-8.8s : %s\n%s", rouge, monNom, where, what, raz)
}
