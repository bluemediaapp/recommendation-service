package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func initialClassifications() {
	app.Post("/initial/:video_id", func(ctx *fiber.Ctx) error {
		// Query data
		videoId, err := strconv.ParseInt(ctx.Params("video_id"), 10, 64)
		if err != nil {
			return err
		}

		// Fetch the db entry
		query := bson.D{{"_id", videoId}}
		rawVideo := videosCollection.FindOne(mctx, query)
		var video databaseVideo
		err = rawVideo.Decode(&video)
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
		return ctx.JSON(result)
	})
}