package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sort"
	"strings"
)

func userClassifications() {
	app.Post("/user/:user_id", func(ctx *fiber.Ctx) error {
		getVideos(databaseUser{Interests: make([]string, 0)})
		return nil
	})
}

func getVideos(user databaseUser) []databaseVideo {
	videos := getAllVideos()
	log.Print(len(videos))
	scoredVideos := make(map[float64]databaseVideo)

	for videoId := range videos {
		video := videos[videoId]
		score := 0.0

		// Interest score
		// This is the average of the tags in your interests list.
		// This should be a multiplier (0-1)

		interestScore := 0.0
		interestsJoined := strings.Join(user.Interests, " ")
		for tagId := range video.Tags {
			tag := video.Tags[tagId]
			if strings.Contains(interestsJoined, tag) {
				interestScore = (interestScore + 1) / 2
			} else {
				interestScore = interestScore / 2
			}
		}

		// Bad topics score
		// A score starting at 1 getting .25 removed for every bad topic
		// Minimum value is .1
		badTopicsScore := 1 - (float64(video.BadTopics) * .25)
		if badTopicsScore > 0 {
			badTopicsScore = .1
		}

		// Calculate score
		score = 1000.0
		score = score * ((interestScore + 1) *.5)
		score = score * badTopicsScore

		// Save score
		scoredVideos[score] = video
	}
	// Sort
	keys := make([]float64, 0, len(scoredVideos))
	for k := range scoredVideos {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	sortedVideos := make([]databaseVideo, len(scoredVideos))
	for score, video := range scoredVideos {
		log.Printf("%s - %f", video.Title, score)
		sortedVideos = append(sortedVideos, video)
	}
	return sortedVideos
}

func getAllVideos() []databaseVideo {
	rawVideos, err := videosCollection.Find(mctx, bson.D{})
	if err != nil {
		log.Printf("find")
		log.Print(err)
	}
	var videos []databaseVideo

	err = rawVideos.All(mctx, &videos)
	if err != nil {
		log.Print(err)
		return nil
	}
	return videos
}
