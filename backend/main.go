package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongoClient *mongo.Client
var pollsCollection *mongo.Collection
var usersCollection *mongo.Collection
var votesCollection *mongo.Collection
var redisClient *redis.Client

func connectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// Ping to actually verify the connection works, not just that the client object was created
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	mongoClient = client
	db := client.Database("polling_app")
	pollsCollection = db.Collection("polls")
	usersCollection = db.Collection("users")
	votesCollection = db.Collection("votes")

	log.Println("Connected to MongoDB")
}

func connectRedis() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Redis URL parse error:", err)
	}

	redisClient = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis ping failed:", err)
	}

	log.Println("Connected to Redis")
}

func main() {
	godotenv.Load()
	connectMongo()
	connectRedis()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "https://live-polling-app-sable.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
	r.POST("/polls", authMiddleware(), createPollHandler)
	r.GET("/polls/:id", getPollHandler)
	r.POST("/polls/:id/vote", voteHandler)
	r.GET("/ws/polls/:id", pollWebSocketHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
