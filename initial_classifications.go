package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func initialClassifications() {
	app.Post("/initial/classify/:video_id", func(ctx *fiber.Ctx) error {
		// Query data
		videoId := ctx.Params("video_id")

		// Fetch the db entry
		query := bson.D{{"_id", videoId}}
		rawVideo := videosCollection.FindOne(mctx, query)
		var video databaseVideo
		err := rawVideo.Decode(&video)
		if err != nil {
			return err
		}
		video.Id = videoId

		// Process the data
		result := initialClassificationResult{}

		// Get the amount of bad tags
		// TODO: Transcriptions?
		for tagId := range video.Tags {
			tag := video.Tags[tagId]
			if isBadTopic(tag) {
				result.BadTopics += 1
			}
		}
		words := strings.Split(video.Description, " ")
		for descriptionId:= range words {
			word := words[descriptionId]
			if isBadTopic(word) {
				result.BadTopics += 1
			}
		}
		return ctx.JSON(result)
	})
}