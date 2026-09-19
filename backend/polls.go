package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type createPollInput struct {
	Question string   `json:"question" binding:"required"`
	Options  []string `json:"options" binding:"required,min=2"`
}

func createPollHandler(c *gin.Context) {
	var input createPollInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")

	poll := Poll{
		Question:  input.Question,
		Options:   input.Options,
		CreatorID: userID,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := pollsCollection.InsertOne(ctx, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create poll"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.InsertedID})
}

func getPollHandler(c *gin.Context) {
	pollID := c.Param("id")

	objID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll Poll
	err = pollsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}

	c.JSON(http.StatusOK, poll)
}
