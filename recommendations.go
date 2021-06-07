package main

import (
	"github.com/bluemediaapp/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

func userClassifications() {
	app.Get("/user/:user_id", func(ctx *fiber.Ctx) error {
		// Get the user
		userId, err := strconv.ParseInt(ctx.Params("user_id"), 10, 64)
		if err != nil {
			return err
		}

		user, err := getUser(userId)
		if err != nil {
			return err
		}
		ignoreRaw := strings.Split(ctx.Query("ignore", ""), "")
		ignores := make([]int64, 0)
		for _, ignoreId := range ignoreRaw {
			id, err := strconv.ParseInt(ignoreId, 10, 64)
			if err != nil {
				return err
			}
			ignores = append(ignores, id)
		}

		videos := getRecommendedVideos(user, ignores)
		return ctx.JSON(videos)
	})
}

func getRecommendedVideos(user models.DatabaseUser, ignore []int64) []models.DatabaseVideo {
	videos := getRandomVideos(100)
	scoredVideos := make(map[float64]models.DatabaseVideo)

	for videoId := range videos {
		video := videos[videoId]

		if hasWatchedVideo(user, video) {
			continue
		}
		if contains(ignore, video.Id) {
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
			calculatedScore := math.Min(float64(thisInterestScore)/100.0, 1)
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
	sortedVideos := make([]models.DatabaseVideo, 0)
	for _, video := range scoredVideos {
		sortedVideos = append(sortedVideos, video)
	}
	if len(sortedVideos) > 10 {
		sortedVideos = sortedVideos[:9]
	}
	return sortedVideos
}

func getRandomVideos(count int) []models.DatabaseVideo {
	// Gets x random videos
	query := []bson.D{bson.D{{"$sample", bson.D{{"size", count}}}}}
	rawVideos, err := videosCollection.Aggregate(mctx, query)
	if err != nil {
		log.Print(err)
		return []models.DatabaseVideo{}
	}
	var videos []models.DatabaseVideo

	err = rawVideos.All(mctx, &videos)
	if err != nil {
		log.Print(err)
		return []models.DatabaseVideo{}
	}
	return videos
}

func hasWatchedVideo(user models.DatabaseUser, video models.DatabaseVideo) bool {
	filter := bson.D{{"user_id", user.Id}, {"video_id", video.Id}}
	var limit int64 = 1
	documentCount, err := watchedVideosCollection.CountDocuments(mctx, filter, &options.CountOptions{
		Limit: &limit,
	})
	if err != nil {
		log.Print(err)
		return true
	}
	return documentCount == int64(1)
}

func getUser(userId int64) (models.DatabaseUser, error) {
	query := bson.D{{"_id", userId}}
	rawVideo := usersCollection.FindOne(mctx, query)
	var user models.DatabaseUser
	err := rawVideo.Decode(&user)
	if err != nil {
		return models.DatabaseUser{}, err
	}
	return user, nil
}

func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
