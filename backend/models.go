package main

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string        `bson:"email" json:"email"`
	PasswordHash string        `bson:"passwordHash" json:"-"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
}

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []string      `bson:"options" json:"options"`
	CreatorID string        `bson:"creatorId" json:"creatorId"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}

type Vote struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	PollID      string        `bson:"pollId" json:"pollId"`
	OptionIndex int           `bson:"optionIndex" json:"optionIndex"`
	VoterID     string        `bson:"voterId" json:"voterId"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
}
