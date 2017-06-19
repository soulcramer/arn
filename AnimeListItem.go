package arn

// AnimeListStatus values for anime list items
const (
	AnimeListStatusWatching  = "watching"
	AnimeListStatusCompleted = "completed"
	AnimeListStatusPlanned   = "planned"
	AnimeListStatusDropped   = "dropped"
	AnimeListStatusHold      = "hold"
)

// AnimeListItem ...
type AnimeListItem struct {
	AnimeID      string      `json:"animeId"`
	Status       string      `json:"status"`
	Episodes     int         `json:"episodes"`
	Rating       AnimeRating `json:"rating"`
	Notes        string      `json:"notes"`
	RewatchCount int         `json:"rewatchCount"`
	Private      bool        `json:"private"`
	Created      string      `json:"created"`
	Edited       string      `json:"edited"`

	anime *Anime
}

// Anime fetches the associated anime data.
func (item *AnimeListItem) Anime() *Anime {
	if item.anime == nil {
		item.anime, _ = GetAnime(item.AnimeID)
	}

	return item.anime
}

// FinalRating returns the overall score for the anime.
func (item *AnimeListItem) FinalRating() float64 {
	return item.Rating.Overall
}
