package utils

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"mime/multipart"

	"github.com/nfnt/resize"
)

func CropAvatarOutOfSkin(skinFileHeader multipart.FileHeader) ([]byte, error) {
	buffer := new(bytes.Buffer)

	skinFile, err := skinFileHeader.Open()
	if err != nil {
		return buffer.Bytes(), err
	}
	defer skinFile.Close()

	skinImage, err := png.Decode(skinFile)
	if err != nil {
		return buffer.Bytes(), err
	}

	faceRect := image.Rect(8, 8, 16, 16)

	face := image.NewRGBA(faceRect)
	draw.Draw(face, faceRect, skinImage, faceRect.Min, draw.Src)

	avatar := resize.Resize(64, 64, face, resize.NearestNeighbor)

	err = png.Encode(buffer, avatar)
	return buffer.Bytes(), err
}
