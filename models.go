package main

type databaseVideo struct {
	// Visible metadata
	Id string `json:"_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedAt int64 `json:"created_at"`
	Series string `json:"series"`
	VideoLength uint8 `json:"video_length"`

	// Backend metadata
	Tags []string `json:"tags"`
	Modifiers []string `json:"modifiers"`
	BadTopics uint8 `json:"bad_topics"`
	Transcription string `json:"transcription"`

}

type initialClassificationResult struct {
	BadTopics uint8 `json:"bad_topics"`
}
