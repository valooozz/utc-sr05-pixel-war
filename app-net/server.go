package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"utils"
)

var ws *websocket.Conn = nil

// INTERFACE -> APP NET
//
//	ne sert pas ici
func doWebsocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	cnx, err := upgrader.Upgrade(w, r, nil)
	ws = cnx
	if err != nil {
		utils.DisplayError("Frontend net", "doWebsocket", "upgrade : "+err.Error())
		return
	}

	for {
		_, message, err := cnx.ReadMessage()
		if err != nil {
			utils.DisplayError("Frontend net", "doWebsocket", "readmessage : "+err.Error())
			return
		}
		traiterMessageInterface(message)
	}
}

// APP NET -> INTERFACE
func wsSend(msg string) {
	if ws != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func launchServer(port string, addr string) {
	http.HandleFunc("/ws", doWebsocket)
	//utils.DisplayError(monNom, "LaunchServer", "Lancement serveur net")
	http.ListenAndServe(addr+":"+port, nil)
}

func traiterMessageInterface(msg []byte) {
	/*message := string(msg)
	if message == "inactif" {
		monEtat = "inactif"
	} else if message == "actif" {
		monEtat = "actif"
	}*/
}
