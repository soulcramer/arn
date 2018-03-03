package imageoutput

import (
	"bytes"
	"errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	"github.com/nfnt/resize"
)

// OriginalFile ...
type OriginalFile struct {
	Directory string
	Size      int
}

// Save writes the original avatar to the file system.
func (output *OriginalFile) Save(avatar *MetaImage, baseName string) error {
	// Determine file extension
	extension := avatar.Extension()

	if extension == "" {
		return errors.New("Unknown format: " + avatar.Format)
	}

	// Resize if needed
	data := avatar.Data
	img := avatar.Image

	if img.Bounds().Dx() > output.Size {
		img = resize.Resize(uint(output.Size), 0, img, resize.Lanczos3)
		buffer := new(bytes.Buffer)

		var err error
		switch extension {
		case ".jpg":
			err = jpeg.Encode(buffer, img, nil)
		case ".png":
			err = png.Encode(buffer, img)
		case ".gif":
			err = gif.Encode(buffer, img, nil)
		}

		if err != nil {
			return err
		}

		data = buffer.Bytes()
	}

	// Write to file
	fileName := output.Directory + baseName + extension
	return ioutil.WriteFile(fileName, data, 0644)
}
