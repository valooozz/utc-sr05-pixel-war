package utils

// Définition des types
const sepM = "/" //séparateur dans les messages
const sepP = "=" //séparateur ddans les paires clé/valeur

type Couleur bool

const (
	Blanc Couleur = false
	Jaune Couleur = true
)

type MessagePixel struct {
	PositionX int
	PositionY int
	Rouge     int
	Vert      int
	Bleu      int
}

type Message struct {
	Pixel       MessagePixel
	Horloge     int
	Vectorielle HorlogeVectorielle
	Nom         string
	Couleur     Couleur
	Prepost     bool // false pour les messages normaux
}

type EtatLocal struct {
	NomSite          string
	Vectorielle      HorlogeVectorielle
	ListMessagePixel []MessagePixel
}

type EtatGlobal struct {
	ListEtatLocal      []EtatLocal
	ListMessagePrepost []Message
}

type MessageEtat struct {
	EtatLocal EtatLocal
	Bilan     int
}

type HorlogeVectorielle map[string]int
