package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	app    = fiber.New()
	config *Config
	client *mongo.Client
	mctx   = context.Background()

	videosCollection        *mongo.Collection
	badTopicsCollection     *mongo.Collection
	watchedVideosCollection *mongo.Collection
)

func main() {
	config = &Config{
		port:     os.Getenv("port"),
		mongoUri: os.Getenv("mongo_uri"),
	}

	initDb()

	// Modules
	initialClassifications()
	userClassifications()

	log.Printf("Running on port 3000!")
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
	badTopicsCollection = db.Collection("bad_topics")
	watchedVideosCollection = db.Collection("watched_videos")
}

func isBadTopic(topic string) bool {
	documentCount, err := badTopicsCollection.CountDocuments(mctx, bson.D{{"topic", topic}})
	if err != nil {
		log.Print(err)
		return false
	}
	return documentCount != 0
}
