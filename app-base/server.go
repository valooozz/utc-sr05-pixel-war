package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"utils"
)

var ws *websocket.Conn = nil

// INTERFACE -> APP BASE
func doWebsocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	cnx, err := upgrader.Upgrade(w, r, nil)
	ws = cnx
	if err != nil {
		utils.DisplayError("Frontend", "doWebsocket", "upgrade : "+err.Error())
		return
	}

	for {
		_, message, err := cnx.ReadMessage()
		if err != nil {
			utils.DisplayError("Frontend", "doWebsocket", "readmessage : "+err.Error())
			return
		}
		traiterMessageInterface(message)
	}
}

// APP BASE -> INTERFACE
func wsSend(msg string) {
	if ws == nil {
		utils.DisplayError("Frontend", "WsSend", "WebSocket non ouverte")
	} else {
		err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			utils.DisplayError("Frontend", "WsSend", "Message non envoyé "+err.Error())
		}
	}
}

func launchServer(port string, addr string) {
	http.HandleFunc("/ws", doWebsocket)
	utils.DisplayError(monNom, "LaunchServer", "Lancement serveur")
	http.ListenAndServe(addr+":"+port, nil)
}

func traiterMessageInterface(msg []byte) {
	message := string(msg)
	if message == "sauvegarder" {
		mutex.Lock()
		envoiSequentiel("sauvegarde")
		mutex.Unlock()
	} else {
		demandeSC()
		wsSend(message)
		envoyerMessage(message)
		relacherSC()
	}
}
