package arn

import (
	"bytes"
	"fmt"
	"image"

	"github.com/animenotifier/arn/imageoutput"
)

// UserAvatar ...
type UserAvatar struct {
	Extension string `json:"extension"`
	Source    string `json:"source"`
}

// RefreshAvatar ...
func (user *User) RefreshAvatar() {
	// TODO: ...
}

// SetAvatarBytes accepts a byte buffer that represents an image file and updates the avatar.
func (user *User) SetAvatarBytes(data []byte) error {
	// Decode
	img, format, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return err
	}

	return user.SetAvatar(&imageoutput.MetaImage{
		Image:  img,
		Format: format,
		Data:   data,
	})
}

// SetAvatar ...
func (user *User) SetAvatar(avatar *imageoutput.MetaImage) error {
	fmt.Println(user.Nick, "uploaded a new avatar:", len(avatar.Data), avatar)

	var lastError error

	// Save the different image formats and sizes
	for _, output := range avatarOutputs {
		err := output.Save(avatar, user.ID)

		if err != nil {
			lastError = err
		}
	}

	user.Avatar.Extension = avatar.Extension()
	return lastError
}
