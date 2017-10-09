package arn

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn/autocorrect"
	"github.com/parnurzeal/gorequest"
)

var youtubeIDRegex = regexp.MustCompile(`youtu(?:.*\/v\/|.*v=|\.be\/)([A-Za-z0-9_-]{11})`)

// SoundCloudToSoundTrack ...
type SoundCloudToSoundTrack struct {
	ID           string `json:"id"`
	SoundTrackID string `json:"soundTrackId"`
}

// YoutubeToSoundTrack ...
type YoutubeToSoundTrack SoundCloudToSoundTrack

// Authorize returns an error if the given API POST request is not authorized.
func (soundtrack *SoundTrack) Authorize(ctx *aero.Context) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Update updates the soundtrack object.
func (soundtrack *SoundTrack) Update(ctx *aero.Context, data interface{}) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	soundtrack.Edited = DateTimeUTC()
	soundtrack.EditedBy = user.ID

	updates := data.(map[string]interface{})
	return SetObjectProperties(soundtrack, updates, nil)
}

// GetSoundCloudMedia returns an ExternalMedia object for the given Soundcloud link.
func GetSoundCloudMedia(url string) (*ExternalMedia, error) {
	var err error
	_, body, errs := gorequest.New().Get("https://api.soundcloud.com/resolve.json?url=" + url + "&client_id=" + APIKeys.SoundCloud.ID).EndBytes()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	var soundcloud SoundCloudTrack
	err = json.Unmarshal(body, &soundcloud)

	if err != nil {
		return nil, err
	}

	if soundcloud.ID == 0 {
		return nil, errors.New("Invalid Soundcloud response as the ID is not valid")
	}

	soundCloudID := strconv.Itoa(soundcloud.ID)

	return &ExternalMedia{
		Service:   "SoundCloud",
		ServiceID: soundCloudID,
	}, nil
}

// GetYoutubeMedia returns an ExternalMedia object for the given Youtube link.
func GetYoutubeMedia(url string) (*ExternalMedia, error) {
	matches := youtubeIDRegex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return nil, errors.New("Invalid Youtube URL")
	}

	videoID := matches[1]

	media := &ExternalMedia{
		Service:   "Youtube",
		ServiceID: videoID,
	}

	return media, nil
}

// Create sets the data for a new soundtrack with data we received from the API request.
func (soundtrack *SoundTrack) Create(ctx *aero.Context) error {
	data, err := ctx.RequestBodyJSONObject()

	if err != nil {
		return err
	}

	userID, ok := ctx.Session().Get("userId").(string)

	if !ok || userID == "" {
		return errors.New("Not logged in")
	}

	user, err := GetUser(userID)

	if err != nil {
		return err
	}

	soundtrack.ID = GenerateID("SoundTrack")
	soundtrack.Likes = []string{}
	soundtrack.Created = DateTimeUTC()
	soundtrack.CreatedBy = user.ID
	soundtrack.Media = []*ExternalMedia{}

	// Soundcloud
	var soundcloud *ExternalMedia
	url, _ := data["soundcloud"].(string)

	if url != "" {
		soundcloud, err = GetSoundCloudMedia(url)

		if err != nil {
			return err
		}

		// Check that the track hasn't been posted yet
		_, err = DB.Get("SoundCloudToSoundTrack", soundcloud.ServiceID)

		if err == nil {
			return errors.New("This Soundcloud track has already been posted")
		}

		// Add to media
		soundtrack.Media = append(soundtrack.Media, soundcloud)
	}

	// Youtube
	var youtube *ExternalMedia
	url, _ = data["youtube"].(string)

	if url != "" {
		youtube, err = GetYoutubeMedia(url)

		if err != nil {
			return err
		}

		// Check that the video hasn't been posted yet
		_, err = DB.Get("YoutubeToSoundTrack", youtube.ServiceID)

		if err == nil {
			return errors.New("This Youtube video has already been posted")
		}

		// Add to media
		soundtrack.Media = append(soundtrack.Media, youtube)
	}

	// Tags
	tags, _ := data["tags"].([]interface{})
	soundtrack.Tags = make([]string, 0)

	animeFound := false
	for i := range tags {
		tag := tags[i].(string)
		tag = autocorrect.FixTag(tag)

		if strings.HasPrefix(tag, "anime:") {
			animeID := strings.TrimPrefix(tag, "anime:")
			_, err := GetAnime(animeID)

			if err != nil {
				return errors.New("Invalid anime ID")
			}

			animeFound = true
		}

		if tag != "" {
			soundtrack.Tags = append(soundtrack.Tags, tag)
		}
	}

	// No media added
	if len(soundtrack.Media) == 0 {
		return errors.New("No media specified (at least 1 media source is required)")
	}

	// No anime found
	if !animeFound {
		return errors.New("Need to specify at least one anime")
	}

	// No tags
	if len(tags) < 1 {
		return errors.New("Need to specify at least one tag")
	}

	// Save Soundcloud reference
	if soundcloud != nil {
		err = DB.Set("SoundCloudToSoundTrack", soundcloud.ServiceID, &SoundCloudToSoundTrack{
			ID:           soundcloud.ServiceID,
			SoundTrackID: soundtrack.ID,
		})

		if err != nil {
			return err
		}
	}

	// Save Youtube reference
	if youtube != nil {
		err = DB.Set("YoutubeToSoundTrack", youtube.ServiceID, &YoutubeToSoundTrack{
			ID:           youtube.ServiceID,
			SoundTrackID: soundtrack.ID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Save saves the soundtrack object in the database.
func (soundtrack *SoundTrack) Save() error {
	return DB.Set("SoundTrack", soundtrack.ID, soundtrack)
}
