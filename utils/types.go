package utils

// Définition des types
const sepM = "/" // séparateur dans les messages
const sepP = "=" // séparateur dans les paires clé/valeur

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

type MessageSauvegarde struct {
	ListMessagePixel []MessagePixel
	Vectorielle      HorlogeVectorielle
}

type Message struct {
	Pixel       MessagePixel
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
}

type HorlogeVectorielle map[string]int

/////////////////////
//Exclusion mutuelle
/////////////////////

type Estampille struct {
	Site    int // numéro du site
	Horloge int // horloge entière
}

// Type de demande d'accès à la section critique (accès, libération)
type TypeSC int

const (
	Requete    TypeSC = 0
	Liberation TypeSC = 1
	Accuse     TypeSC = 2
)

// Message pour la demande d'accès à la section critique
type MessageExclusionMutuelle struct {
	Type       TypeSC
	Estampille Estampille
}

type MessageAccuse struct {
	SiteCible  int
	Estampille Estampille
}

/////////////////////
// Messages net
/////////////////////

type Header struct {
	Origine     int
	Destination int
	Initiateur  int
	Vecteur     []int
}

type MessageNet struct {
	Header         Header
	MessageControl string
}

/////////////////////
// Messages id
/////////////////////

type MessageId struct {
	Id      int
	Message string
}

/////////////////////
// Routage
/////////////////////

type Route struct {
	Origine     int
	Destination int
}

type TableDeRoutage []Route

////////////////
// Election
////////////////

type ColorationVague int

const (
	Bleu  ColorationVague = 1
	Rouge ColorationVague = 2
	Vert  ColorationVague = 3
)

type MessageVague struct {
	Site          int             // Le site qui envoie le MessageVague
	Coloration    ColorationVague // Le type de message (bleu, rouge, vert)
	Info          int             // le numéro du potentiel élu pour les messages bleus et rouges / le nouveau N pour les messages verts
	Cible         int             // le site cible des messages rouges / le site qui doit ignorer (parent) les messages bleus et verts
	SiteDemandeur int             // le site qui fait la demande de raccordement
}

/////////////////
// Raccordement
/////////////////

type MessageRaccord struct {
	Site  int    // Le site qui envoie le MessageRaccord
	Type  string // demande ou acceptation
	Info  int    // 1 si on veut rejoindre / -1 si on veut quitter
	Cible int    // Le site à qui est destiné le MessageRaccord
}

type Demande struct {
	Site int // Le site qui a fait la demande
	Info int // 1 si c'est pour rejoindre / -1 si c'est pour quitter
}
