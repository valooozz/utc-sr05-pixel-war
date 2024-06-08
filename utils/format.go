package utils

import (
	"fmt"
	"strconv"
	"strings"
)

///////////////
// MessagePixel
///////////////

func MessagePixelToString(pixel MessagePixel) string {
	return sepM + sepP + "positionX" + sepP + strconv.Itoa(pixel.PositionX) + sepM + sepP + "positionY" + sepP +
		strconv.Itoa(pixel.PositionY) + sepM + sepP + "R" + sepP + strconv.Itoa(pixel.Rouge) + sepM + sepP + "G" +
		sepP + strconv.Itoa(pixel.Vert) + sepM + sepP + "B" + sepP + strconv.Itoa(pixel.Bleu)
}

func StringToMessagePixel(str string) MessagePixel {
	posX, _ := strconv.Atoi(TrouverValeur(str, "positionX"))
	posY, _ := strconv.Atoi(TrouverValeur(str, "positionY"))
	r, _ := strconv.Atoi(TrouverValeur(str, "R"))
	v, _ := strconv.Atoi(TrouverValeur(str, "G"))
	b, _ := strconv.Atoi(TrouverValeur(str, "B"))

	messagepixel := MessagePixel{posX, posY, r, v, b}
	return messagepixel
}

////////////
// MessageSauvegarde
////////////

func MessageSauvegardeToString(sauvegarde MessageSauvegarde) string {
	sep1 := "#"
	sep2 := ";"
	l := ""
	for _, messagePixel := range sauvegarde.ListMessagePixel {
		l += "_"
		l += MessagePixelToString(messagePixel)
	}

	return sep1 + sep2 + "vectorielle" + sep2 + HorlogeVectorielleToString(sauvegarde.Vectorielle) + sep1 + sep2 + "listSauvegarde" + sep2 + l
}

func StringToMessageSauvegarde(str string) MessageSauvegarde {
	listSauvegarde := TrouverValeur(str, "listSauvegarde")
	tabListSauvegarde := strings.Split(listSauvegarde, "_")

	messageSauvegarde := MessageSauvegarde{}
	var liste []MessagePixel

	for _, strMessagePixel := range tabListSauvegarde {
		if strMessagePixel != "" {
			liste = append(liste, StringToMessagePixel(strMessagePixel))
		}
	}

	messageSauvegarde.ListMessagePixel = liste
	messageSauvegarde.Vectorielle = StringToHorlogeVectorielle(TrouverValeur(str, "vectorielle"))

	return messageSauvegarde
}

//////////
// Message
//////////

func MessageToString(message Message) string {
	c := ""
	if message.Couleur {
		c = "jaune"
	} else {
		c = "blanc"
	}
	return MessagePixelToString(message.Pixel) + sepM + sepP + "vectorielle" + sepP + HorlogeVectorielleToString(message.Vectorielle) + sepM + sepP + "nom" + sepP + message.Nom + sepM + sepP + "couleur" + sepP + c +
		sepM + sepP + "prepost" + sepP + strconv.FormatBool(message.Prepost)

}

func StringToMessage(str string) Message {
	messagepixel := StringToMessagePixel(str)
	hv := TrouverValeur(str, "vectorielle")
	n := TrouverValeur(str, "nom")
	cV := TrouverValeur(str, "couleur")
	var c Couleur
	if cV == "jaune" {
		c = Jaune
	} else {
		c = Blanc
	}
	prep, _ := strconv.ParseBool(TrouverValeur(str, "prepost"))
	message := Message{messagepixel, StringToHorlogeVectorielle(hv), n, c, prep}
	return message
}

//////////////
// MessageEtat
//////////////

func MessageEtatToString(etat MessageEtat) string {
	sep1 := "~"
	sep2 := ","
	return sep1 + sep2 + "etat" + sep2 + EtatLocalToString(etat.EtatLocal)
}

func StringToMessageEtat(str string) MessageEtat {
	etatLocal := StringToEtatLocal(TrouverValeur(str, "etat"))

	return MessageEtat{etatLocal}
}

////////////
// EtatLocal
////////////

func EtatLocalToString(etatLocal EtatLocal) string {
	sep1 := "#"
	sep2 := ";"
	l := ""
	for _, messagePixel := range etatLocal.ListMessagePixel {
		l += "_"
		l += MessagePixelToString(messagePixel)
	}

	return sep1 + sep2 + "nom" + sep2 + etatLocal.NomSite +
		sep1 + sep2 + "vectorielle" + sep2 + HorlogeVectorielleToString(etatLocal.Vectorielle) +
		sep1 + sep2 + "liste" + sep2 + l
}

func StringToEtatLocal(str string) EtatLocal {
	var liste []MessagePixel
	listeMessagePixel := TrouverValeur(str, "liste")
	strVectorielle := TrouverValeur(str, "vectorielle")
	tabListeMessagePixel := strings.Split(listeMessagePixel, "_")

	for _, strMessagePixel := range tabListeMessagePixel {
		if strMessagePixel != "" {
			liste = append(liste, StringToMessagePixel(strMessagePixel))
		}
	}

	return EtatLocal{TrouverValeur(str, "nom"), StringToHorlogeVectorielle(strVectorielle), liste}
}

/////////////////////
// HorlogeVectorielle
/////////////////////

func HorlogeVectorielleToString(horloge HorlogeVectorielle) string {
	sep1 := "_"
	sep2 := ":"
	str := ""

	for site := range horloge {
		str += sep1
		str += site
		str += sep2
		str += strconv.Itoa(horloge[site])
	}

	return str
}

func StringToHorlogeVectorielle(str string) HorlogeVectorielle {
	horloge := HorlogeVectorielle{}
	listeSites := strings.Split(str, "_")

	for _, strSite := range listeSites {
		if strSite != "" {
			hSite := strings.Split(strSite, ":")
			horloge[hSite[0]], _ = strconv.Atoi(hSite[1])
		}
	}

	return horloge
}

///////////////////////////
// MessageExclusionMutuelle
///////////////////////////

func MessageExclusionMutuelleToString(exclumutuelle MessageExclusionMutuelle) string {
	return sepM + sepP + "typeSC" + sepP + strconv.Itoa(int(exclumutuelle.Type)) + sepM + sepP + "estampilleSite" + sepP +
		strconv.Itoa(exclumutuelle.Estampille.Site) + sepM + sepP + "estampilleHorloge" +
		sepP + strconv.Itoa(exclumutuelle.Estampille.Horloge)
}

func StringToMessageExclusionMutuelle(str string) MessageExclusionMutuelle {
	typeSC, _ := strconv.Atoi(TrouverValeur(str, "typeSC"))
	site, _ := strconv.Atoi(TrouverValeur(str, "estampilleSite"))
	horloge, _ := strconv.Atoi(TrouverValeur(str, "estampilleHorloge"))
	estampille := Estampille{site, horloge}
	messageecxlumutuelle := MessageExclusionMutuelle{TypeSC(typeSC), estampille}
	return messageecxlumutuelle
}

////////////////
// MessageTypeSC
////////////////

func MessageTypeSCToString(exclumutuelle TypeSC) string {
	return sepM + sepP + "typeSC" + sepP + strconv.Itoa(int(exclumutuelle))
}

func StringToMessageTypeSC(str string) TypeSC {
	t, _ := strconv.Atoi(TrouverValeur(str, "typeSC"))
	messageecxlumutuelle := TypeSC(t)
	return messageecxlumutuelle
}

////////////////
// MessageAccuse
////////////////

func MessageAccuseToString(message MessageAccuse) string {
	return sepM + sepP + "siteCible" + sepP + strconv.Itoa(message.SiteCible) + sepM + sepP + "estampilleSite" + sepP +
		strconv.Itoa(message.Estampille.Site) + sepM + sepP + "estampilleHorloge" +
		sepP + strconv.Itoa(message.Estampille.Horloge)
}

func StringToMessageAccuse(str string) MessageAccuse {
	s, _ := strconv.Atoi(TrouverValeur(str, "siteCible"))
	site, _ := strconv.Atoi(TrouverValeur(str, "estampilleSite"))
	horloge, _ := strconv.Atoi(TrouverValeur(str, "estampilleHorloge"))
	estampille := Estampille{site, horloge}
	message := MessageAccuse{s, estampille}
	return message
}

////////////////
// MessageNet
////////////////

func VecteurToString(vect []int) string {
	var str string
	str += "["
	for i, val := range vect {
		if i > 0 {
			str += "_"
		}
		str += strconv.Itoa(val)
	}
	str += "]"
	return str
}

func StringToVecteur(s string) []int {
	s = strings.Trim(s, "[]")
	valuesStr := strings.Split(s, "_")
	var vect []int
	for _, valStr := range valuesStr {
		val, _ := strconv.Atoi(valStr)
		vect = append(vect, val)
	}
	return vect
}

func HeaderToString(header Header) string {
	sep1 := "$"
	sep2 := "^"
	return sep1 + sep2 + "origine" + sep2 + strconv.Itoa(header.Origine) + sep1 + sep2 + "destination" + sep2 + strconv.Itoa(header.Destination) +
		sep1 + sep2 + "initiateur" + sep2 + strconv.Itoa(header.Initiateur) + sep1 + sep2 + "vecteur" + sep2 + VecteurToString(header.Vecteur)
}

func StringToHeader(str string) Header {
	o, _ := strconv.Atoi(TrouverValeur(str, "origine"))
	d, _ := strconv.Atoi(TrouverValeur(str, "destination"))
	i, _ := strconv.Atoi(TrouverValeur(str, "initiateur"))
	v := StringToVecteur(TrouverValeur(str, "vecteur"))
	header := Header{o, d, i, v}
	return header
}

func MessageNetToString(message MessageNet) string {
	sep1 := "@"
	sep2 := "+"
	return sep1 + sep2 + "header" + sep2 + HeaderToString(message.Header) + sep1 + sep2 + "messageControl" + sep2 +
		message.MessageControl
}

func StringToMessageNet(str string) MessageNet {
	header := StringToHeader(TrouverValeur(str, "header"))
	messageControl := TrouverValeur(str, "messageControl")
	message := MessageNet{header, messageControl}
	return message
}

////////////////
// MessageId
////////////////

func MessageIdToString(messageId MessageId) string {
	sep1 := "$"
	sep2 := "^"
	return sep1 + sep2 + "id" + sep2 + strconv.Itoa(messageId.Id) + sep1 + sep2 + "message" + sep2 + messageId.Message
}

func StringToMessageId(str string) MessageId {
	id, _ := strconv.Atoi(TrouverValeur(str, "id"))
	message := TrouverValeur(str, "message")
	messageId := MessageId{Id: id, Message: message}
	return messageId
}

/////////////////////
// Routage
/////////////////////

func StringToTableDeRoutage(s string) TableDeRoutage {
	s = strings.Trim(s, "[]") // Enlever les crochets
	routesStr := strings.Split(s, ";")
	var tdr TableDeRoutage
	for _, routeStr := range routesStr {
		fields := strings.Split(routeStr, ",")
		origine, _ := strconv.Atoi(fields[0])
		destination, _ := strconv.Atoi(fields[1])
		tdr = append(tdr, Route{Origine: origine, Destination: destination})
	}
	return tdr
}

func TableDeRoutageToString(tdr TableDeRoutage) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, route := range tdr {
		if i > 0 {
			sb.WriteString(";")
		}
		sb.WriteString(fmt.Sprintf("%d,%d", route.Origine, route.Destination))
	}
	sb.WriteString("]")
	return sb.String()
}

////////////
// Election
////////////

func MessageVagueToString(messageVague MessageVague) string {
	return sepM + sepP + "site" + sepP + strconv.Itoa(messageVague.Site) +
		sepM + sepP + "coloration" + sepP + strconv.Itoa(int(messageVague.Coloration)) +
		sepM + sepP + "info" + sepP + strconv.Itoa(messageVague.Info) +
		sepM + sepP + "cible" + sepP + strconv.Itoa(messageVague.Cible)
}

func StringToMessageVague(str string) MessageVague {
	site, _ := strconv.Atoi(TrouverValeur(str, "site"))
	coloration, _ := strconv.Atoi(TrouverValeur(str, "coloration"))
	info, _ := strconv.Atoi(TrouverValeur(str, "info"))
	cible, _ := strconv.Atoi(TrouverValeur(str, "cible"))

	messageVague := MessageVague{site, ColorationVague(coloration), info, cible}
	return messageVague
}

//////////////////
// Raccordement
//////////////////

func MessageRaccordToString(messageRaccord MessageRaccord) string {
	return sepM + sepP + "site" + sepP + strconv.Itoa(messageRaccord.Site) +
		sepM + sepP + "type" + sepP + messageRaccord.Type +
		sepM + sepP + "info" + sepP + strconv.Itoa(messageRaccord.Info) +
		sepM + sepP + "cible" + sepP + strconv.Itoa(messageRaccord.Cible)
}

func StringToMessageRaccord(str string) MessageRaccord {
	site, _ := strconv.Atoi(TrouverValeur(str, "site"))
	typeM := TrouverValeur(str, "type")
	info, _ := strconv.Atoi(TrouverValeur(str, "info"))
	cible, _ := strconv.Atoi(TrouverValeur(str, "cible"))

	messageRaccord := MessageRaccord{site, typeM, info, cible}
	return messageRaccord
}

func MessageBlocageToString(messageBlocage MessageBlocage) string {
	var c string
	if messageBlocage.Blocage {
		c = "noir"
	} else {
		c = "gris"
	}

	return sepM + sepP + "site" + sepP + strconv.Itoa(messageBlocage.Site) +
		sepM + sepP + "blocage" + sepP + c +
		sepM + sepP + "cible" + sepP + strconv.Itoa(messageBlocage.Cible)
}

func StringToMessageBlocage(str string) MessageBlocage {
	site, _ := strconv.Atoi(TrouverValeur(str, "site"))
	blocage := TrouverValeur(str, "blocage")
	cible, _ := strconv.Atoi(TrouverValeur(str, "cible"))

	var b CouleurBlocage
	if blocage == "noir" {
		b = Noir
	} else {
		b = Gris
	}

	messageBlocage := MessageBlocage{site, b, cible}
	return messageBlocage
}

func TabSCToString(tab []MessageExclusionMutuelle) string {
	var sb = "["

	for i, msg := range tab {
		if i > 0 {
			sb += "|"
		}
		// Conversion de chaque message en chaîne et ajout au StringBuilder
		sb += "T:" + strconv.Itoa(int(msg.Type)) + "S:" + strconv.Itoa(msg.Estampille.Site) + "H:" + strconv.Itoa(msg.Estampille.Horloge)
	}
	sb += "]"

	return sb
}
