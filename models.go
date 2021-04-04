package main

// Videos
type databaseVideo struct {
	// Visible metadata
	Id          int64 `bson:"_id" json:"id"`
	Description string `bson:"description" json:"description"`
	Series      string `bson:"series" json:"series"`
	VideoLength int32  `bson:"video_length" json:"video_length"`
	Public      bool   `bson:"public" json:"public"`

	// Backend metadata
	Tags          []string `bson:"tags" json:"tags"`
	Modifiers     []string `bson:"modifiers" json:"modifiers"`
	BadTopics     int32    `bson:"bad_topics" json:"bad_topics"`
}

type initialClassificationResult struct {
	BadTopics uint8 `bson:"bad_topics" json:"bad_topics"`
}

// Users
type databaseUser struct {
	// Visible metadata
	Id       int64 `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`

	// Backend metadata
	Interests []string `bson:"interests" json:"interests"`
}

// Watching
type databaseWatchEvent struct {
	VideoId     int64 `bson:"video_id" json:"video_id"`
	UserId      int64 `bson:"user_id" json:"user_id"`
	TimeWatched int32  `bson:"time_watched" json:"time_watched"`
}

type databaseLikeEvent struct {
	VideoId int64 `bson:"video_id" json:"video_id"`
	UserId  int64 `bson:"user_id" json:"user_id"`
}
