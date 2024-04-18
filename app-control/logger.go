package main

import (
	"log"
	"os"
)

//Définition des fonctions de différentes displays

var stderr = log.New(os.Stderr, "", 0)

func displayInfo(where string, what string) {
	stderr.Printf("%s + %-8.8s : %s\n", nom, where, what)
}

func displayWarning(where string, what string) {

	stderr.Printf("%s %s * %-8.8s : %s\n%s", orange, nom, where, what, raz)
}

func displayError(where string, what string) {
	stderr.Printf("%s %s ! %-8.8s : %s\n%s", rouge, nom, where, what, raz)
}
