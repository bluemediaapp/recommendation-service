package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	app     = fiber.New()
	config  *Config
	client  *mongo.Client
	mctx    = context.Background()

	videosCollection        *mongo.Collection
	watchedVideosCollection *mongo.Collection
	usersCollection         *mongo.Collection
)

func main() {
	config = &Config{
		port:     os.Getenv("port"),
		mongoUri: os.Getenv("mongo_uri"),
	}

	rand.Seed(time.Now().UnixNano())

	initDb()

	// Modules
	userClassifications()

	log.Fatal(app.Listen(config.port))
}

func initDb() {
	// Connect mongo
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.mongoUri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(mctx)
	if err != nil {
		log.Fatal(err)
	}

	// Setup tables
	db := client.Database("blue")
	videosCollection = db.Collection("video_metadata")
	watchedVideosCollection = db.Collection("watched_videos")
	usersCollection = db.Collection("users")
}
