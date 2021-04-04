package main

// Videos
type databaseVideo struct {
	// Visible metadata
	Id          string `json:"_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	Series      string `json:"series"`
	VideoLength int32  `json:"video_length"`
	Public      bool   `json:"public"`

	// Backend metadata
	Tags          []string `json:"tags"`
	Modifiers     []string `json:"modifiers"`
	BadTopics     int32    `json:"bad_topics"`
	Transcription string   `json:"transcription"`
}

type initialClassificationResult struct {
	BadTopics uint8 `json:"bad_topics"`
}

// Users
type databaseUser struct {
	// Visible metadata
	Id       string `json:"_id"`
	Username string `json:"username"`

	// Backend metadata
	CreatedAt int64    `json:"created_at"`
	Interests []string `json:"interests"`
}

// Watching
type databaseWatchEvent struct {
	VideoId     string `json:"video_id"`
	UserId      string `json:"user_id"`
	TimeWatched uint8  `json:"time_watched"`
}

type databaseLikeEvent struct {
	VideoId string `json:"video_id"`
	UserId  string `json:"user_id"`
}
