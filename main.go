package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func broadcastMessage(sender *websocket.Conn, msgType int, msg []byte) {

	for client := range clients {
		if client != sender {
			err := client.WriteMessage(msgType, msg)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ws", func(ctx *gin.Context) {
		conn, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Println("Error while upgrading:", err)
			return
		}
		fmt.Println(conn)
		clients[conn] = true
		fmt.Println("New client connected")

		defer func() {
			delete(clients, conn)
			conn.Close()
		}()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error while reading message:", err)
				break
			}

			fmt.Println("Message Received: ", string(msg))

			if string(msg) == "ping" {
				err := conn.WriteMessage(msgType, []byte("pong"))
				if err != nil {
					fmt.Println("Error while writing pong:", err)
					break
				}
			} else {
				broadcastMessage(conn, msgType, msg)
			}
		}
	})

	fmt.Println("Server started on port 8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error while starting server:", err)
	}
}
