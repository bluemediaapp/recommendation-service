package main

// Videos
type databaseVideo struct {
	// Visible metadata
	Id          string `bson:"_id" json:"id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	CreatedAt   int64  `bson:"created_at" json:"created_at"`
	Series      string `bson:"series" json:"series"`
	VideoLength int32  `bson:"video_length" json:"video_length"`
	Public      bool   `bson:"public" json:"public"`

	// Backend metadata
	Tags          []string `bson:"tags" json:"tags"`
	Modifiers     []string `bson:"modifiers" json:"modifiers"`
	BadTopics     int32    `bson:"bad_topics" json:"bad_topics"`
	Transcription string   `bson:"transcription" json:"transcription"`
}

type initialClassificationResult struct {
	BadTopics uint8 `bson:"bad_topics" json:"bad_topics"`
}

// Users
type databaseUser struct {
	// Visible metadata
	Id       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`

	// Backend metadata
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	Interests []string `bson:"interests" json:"interests"`
}

// Watching
type databaseWatchEvent struct {
	VideoId     string `bson:"video_id" json:"video_id"`
	UserId      string `bson:"user_id" json:"user_id"`
	TimeWatched uint8  `bson:"time_watched" json:"time_watched"`
}

type databaseLikeEvent struct {
	VideoId string `bson:"video_id" json:"video_id"`
	UserId  string `bson:"user_id" json:"user_id"`
}
