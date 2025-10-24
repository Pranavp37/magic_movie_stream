package models

type Movie struct {
	ID          string        `bson:"_id,omitempty" json:"_id"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id"`
	Title       string        `bson:"title" json:"title"`
	PosterPath  string        `bson:"poster_path" json:"poster_path"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id"`
	Genre       []string      `bson:"genre" json:"genre"`
	AdminReview string        `bson:"admin_review" json:"admin_review"`
	Ranking     RankingDetail `bson:"ranking" json:"ranking"`
}

type RankingDetail struct {
	Score   float64 `bson:"score" json:"score,omitempty"`
	Reviews int     `bson:"reviews" json:"reviews,omitempty"`
}
type MovieListResponse struct {
	Movies []Movie `bson:"movies" json:"movies"`
}			