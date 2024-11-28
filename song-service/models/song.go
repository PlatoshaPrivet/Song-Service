package models

// Song represents a music song
// @description Song model definition
type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group" binding:"required"`
	Song        string `json:"song" binding:"required"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

//For updateSong func
type Required struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
