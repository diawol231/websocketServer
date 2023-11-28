package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

var (
	clients map[*websocket.Conn]bool
	adminConn *websocket.Conn
	arrayInt [2]int
)

func main() {
	clients = make(map[*websocket.Conn]bool)
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		connection, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		defer connection.Close()
		if _, message, _ := connection.ReadMessage(); string(message) == "admin"{
			adminConn = connection
			fmt.Print("Admin connected!\n")
		} else {
			clients[connection] = true
			fmt.Print("Client connected!\n")
		}
		
		defer delete(clients, connection)
		for {
			// Read message from browser
			msgType, msg, err := connection.ReadMessage()
			if err != nil {
				return
			}

			msgList := strings.Split(string(msg), " ")
			if(msgList[0] == "result"){
				adminConn.WriteMessage(msgType, []byte(msgList[1]))
				continue
			}
			// Print the message to the console
			// Write message back to browser
			i := 0
			for conn := range clients {
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msgList[i]))
				conn.WriteMessage(msgType, []byte(msgList[i]))
				i++
				if i >= 2 {
					break
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web.html")
	})

	http.ListenAndServe(":8081", nil)
}
