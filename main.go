package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
	rclient *redis.Client

	videosCollection        *mongo.Collection
	badTopicsCollection     *mongo.Collection
	watchedVideosCollection *mongo.Collection
	usersCollection         *mongo.Collection
)

func main() {
	config = &Config{
		port:     os.Getenv("port"),
		mongoUri: os.Getenv("mongo_uri"),
		redisUri: os.Getenv("redis_uri"),
	}

	rand.Seed(time.Now().UnixNano())

	initDb()
	initRedis()

	// Modules
	initialClassifications()
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
	badTopicsCollection = db.Collection("bad_topics")
	watchedVideosCollection = db.Collection("watched_videos")
	usersCollection = db.Collection("users")
}

func initRedis()  {
	rclient = redis.NewClient(&redis.Options{
		Addr: config.redisUri,
	})
	err := rclient.Ping(mctx).Err()
	if err != nil {
		panic(err)
	}
}

func isBadTopic(topic string) bool {
	documentCount, err := badTopicsCollection.CountDocuments(mctx, bson.D{{"topic", topic}})
	if err != nil {
		log.Print(err)
		return false
	}
	return documentCount != 0
}
