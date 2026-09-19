package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type voteInput struct {
	OptionIndex int `json:"optionIndex" binding:"min=0"`
}

func voteHandler(c *gin.Context) {
	pollID := c.Param("id")

	objID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll id"})
		return
	}

	var input voteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// confirm poll exists and the option index is actually valid for it
	var poll Poll
	if err := pollsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	if input.OptionIndex >= len(poll.Options) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid option index"})
		return
	}

	// record the vote in Mongo (permanent record)
	vote := Vote{
		PollID:      pollID,
		OptionIndex: input.OptionIndex,
		VoterID:     c.ClientIP(), // simple for now — later could be a session/user id
		CreatedAt:   time.Now(),
	}
	if _, err := votesCollection.InsertOne(ctx, vote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not record vote"})
		return
	}

	// atomically bump the Redis counter for this option
	countKey := "poll:" + pollID + ":counts"
	field := strconv.Itoa(input.OptionIndex)
	if err := redisClient.HIncrBy(ctx, countKey, field, 1).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update counts"})
		return
	}

	// fetch full current counts and publish them so all connected clients update
	counts, err := redisClient.HGetAll(ctx, countKey).Result()
	if err != nil {
		log.Println("HGetAll failed:", err)
	} else {
		payload, _ := json.Marshal(counts)
		n, pubErr := redisClient.Publish(ctx, "poll:"+pollID+":updates", payload).Result()
		if pubErr != nil {
			log.Println("Publish failed:", pubErr)
		} else {
			log.Println("Published to", n, "subscribers")
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "vote recorded"})
}
