package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func handlerWebsocket(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()

	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true
	fmt.Println("Client connected")

	for {
		// Read the message from the WebSocket connection
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, ws)
			break
		}

		// Print the message to the console
		fmt.Printf("Message received: %s\n", msg)

		msg = []byte("You said: " + string(msg))

		// Send the message to all clients
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println(err)
				delete(clients, client)
				break
			}
		}
	}

	return nil
}

func handlerTest(c echo.Context) error {
	var msg = []byte("Hello, World!")

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			delete(clients, client)
			break
		}
	}

	return c.String(http.StatusOK, "Test message sent")
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.GET("/ws", handlerWebsocket)

	e.GET("/test", handlerTest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
