package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"sort"
	"strconv"
)

func userClassifications() {
	app.Post("/user/:user_id", func(ctx *fiber.Ctx) error {
		// Get the user
		userId, err := strconv.ParseInt(ctx.Params("user_id"), 10, 64)
		if err != nil {
			return err
		}

		user, err := getUser(userId)
		if err != nil {
			return err
		}

		videos := getRecommendedVideos(user)
		return ctx.JSON(videos)
	})
}

func getRecommendedVideos(user databaseUser) []databaseVideo {
	videos := getRandomVideos(100)
	scoredVideos := make(map[float64]databaseVideo)

	for videoId := range videos {
		video := videos[videoId]

		if hasWatchedVideo(user, video) {
			continue
		}

		score := 0.0

		// Interest score
		// This is the average of the tags in your interests list.
		// This should be a multiplier (0-1)

		interestScore := 0.0
		for tagId := range video.Tags {
			tag := video.Tags[tagId]
			thisInterestScore, exists := user.Interests[tag]
			if !exists {
				interestScore = interestScore / 2
				continue
			}
			calculatedScore := math.Min(float64(thisInterestScore) /100.0, 1)
			interestScore = (interestScore + calculatedScore) / 2
		}
		if interestScore < 0 {
			interestScore = 0
		}

		// Bad topics score
		// A score starting at 1 getting .25 removed for every bad topic
		// Minimum value is .1
		badTopicsScore := 1 - (float64(video.BadTopics) * .25)
		if badTopicsScore <= 0 {
			badTopicsScore = .1
		}

		// Calculate score
		score = 1000.0
		score = score * ((interestScore * 10) + .1)
		score = score * badTopicsScore

		// Save score
		scoredVideos[score] = video
	}
	// Sort
	keys := make([]float64, 0)
	for k := range scoredVideos {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	sortedVideos := make([]databaseVideo, 0)
	for _, video := range scoredVideos {
		sortedVideos = append(sortedVideos, video)
	}
	return sortedVideos
}

func getRandomVideos(count int) []databaseVideo {
	// Gets x random videos
	query := []bson.D{bson.D{{"$sample", bson.D{{"size", count}}}}}
	rawVideos, err := videosCollection.Aggregate(mctx, query)
	if err != nil {
		log.Print(err)
		return []databaseVideo{}
	}
	var videos []databaseVideo

	err = rawVideos.All(mctx, &videos)
	if err != nil {
		log.Print(err)
		return []databaseVideo{}
	}
	return videos
}

func hasWatchedVideo(user databaseUser, video databaseVideo) bool {
	filter := bson.D{{"user_id", user.Id}, {"video_id", video.Id}}
	documentCount, err := watchedVideosCollection.CountDocuments(mctx, filter)
	if err != nil {
		log.Print(err)
		return true
	}
	return documentCount == int64(1)
}

func getUser(userId int64) (databaseUser, error) {
	query := bson.D{{"_id", userId}}
	rawVideo := usersCollection.FindOne(mctx, query)
	var user databaseUser
	err := rawVideo.Decode(&user)
	if err != nil {
		return databaseUser{}, err
	}
	return user, nil
}
