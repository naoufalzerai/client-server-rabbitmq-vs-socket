package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func main() {

	mnt := 100
	//selectionner facture
	fac := uuid.NewV4()
	log.Println("Selection de la facture :" + fac.String())
	log.Println("Montant encaissé :" + fmt.Sprint(mnt))

	//creer l'uuid
	piece := uuid.NewV4()
	log.Println("Numero de la piece :" + piece.String())

	// init
	u := url.URL{Scheme: "ws", Host: "localhost:1323", Path: "/ws"}
	log.Printf("connecting au serveur %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	encaissement := map[string]string{
		"NumPiece": piece.String(),
		"NumFac":   fac.String(),
		"Mnt":      fmt.Sprint(mnt),
		"Caissier": "somei",
	}

	jsonEncaissent, err := json.Marshal(encaissement)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// envoi d'un encaissement
	err = c.WriteMessage(websocket.TextMessage, jsonEncaissent)
	if err != nil {
		log.Fatal("dial:", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal("dial:", err)
		}
		log.Println(string(message))
	}

}
