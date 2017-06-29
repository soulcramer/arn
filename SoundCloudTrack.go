package arn

// SoundCloudTrack ...
type SoundCloudTrack struct {
	Kind                string      `json:"kind"`
	ID                  int         `json:"id"`
	CreatedAt           string      `json:"created_at"`
	UserID              int         `json:"user_id"`
	Duration            int         `json:"duration"`
	Commentable         bool        `json:"commentable"`
	State               string      `json:"state"`
	OriginalContentSize int         `json:"original_content_size"`
	LastModified        string      `json:"last_modified"`
	Sharing             string      `json:"sharing"`
	TagList             string      `json:"tag_list"`
	Permalink           string      `json:"permalink"`
	Streamable          bool        `json:"streamable"`
	EmbeddableBy        string      `json:"embeddable_by"`
	PurchaseURL         interface{} `json:"purchase_url"`
	PurchaseTitle       interface{} `json:"purchase_title"`
	LabelID             interface{} `json:"label_id"`
	Genre               string      `json:"genre"`
	Title               string      `json:"title"`
	Description         string      `json:"description"`
	LabelName           string      `json:"label_name"`
	Release             string      `json:"release"`
	TrackType           string      `json:"track_type"`
	KeySignature        string      `json:"key_signature"`
	Isrc                interface{} `json:"isrc"`
	VideoURL            interface{} `json:"video_url"`
	Bpm                 interface{} `json:"bpm"`
	ReleaseYear         interface{} `json:"release_year"`
	ReleaseMonth        interface{} `json:"release_month"`
	ReleaseDay          interface{} `json:"release_day"`
	OriginalFormat      string      `json:"original_format"`
	License             string      `json:"license"`
	URI                 string      `json:"uri"`
	User                struct {
		ID           int    `json:"id"`
		Kind         string `json:"kind"`
		Permalink    string `json:"permalink"`
		Username     string `json:"username"`
		LastModified string `json:"last_modified"`
		URI          string `json:"uri"`
		PermalinkURL string `json:"permalink_url"`
		AvatarURL    string `json:"avatar_url"`
	} `json:"user"`
	PermalinkURL     string `json:"permalink_url"`
	ArtworkURL       string `json:"artwork_url"`
	StreamURL        string `json:"stream_url"`
	DownloadURL      string `json:"download_url"`
	PlaybackCount    int    `json:"playback_count"`
	DownloadCount    int    `json:"download_count"`
	FavoritingsCount int    `json:"favoritings_count"`
	RepostsCount     int    `json:"reposts_count"`
	CommentCount     int    `json:"comment_count"`
	Downloadable     bool   `json:"downloadable"`
	WaveformURL      string `json:"waveform_url"`
	AttachmentsURI   string `json:"attachments_uri"`
}