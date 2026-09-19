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

	ctx := c.Request.Context()
	channel := "poll:" + pollID + ":updates"
	sub := redisClient.Subscribe(ctx, channel)
	defer sub.Close()

	ch := sub.Channel()

	log.Println("Subscribed to", channel)
	for msg := range ch {
		log.Println("Received on", channel, ":", msg.Payload)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Println("WebSocket write failed, closing:", err)
			return
		}
	}
	log.Println("Subscription loop ended for", channel)
}
