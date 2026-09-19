package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // fine for now, tighten before submission
}

func pollWebSocketHandler(c *gin.Context) {
	pollID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	channel := "poll:" + pollID + ":updates"
	sub := redisClient.Subscribe(c.Request.Context(), channel)
	defer sub.Close()

	ch := sub.Channel()

	// detect disconnects: gorilla requires an active reader to notice a closed connection
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return // client disconnected or errored
			}
		}
	}()

	log.Println("Subscribed to", channel)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Println("WebSocket write failed, closing:", err)
				return
			}
		case <-done:
			log.Println("Client disconnected, closing subscription for", channel)
			return
		}
	}
}
