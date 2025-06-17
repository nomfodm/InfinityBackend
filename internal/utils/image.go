package utils

import (
	"bytes"
	"fmt"
	"github.com/mineatar-io/skin-render"
	"image"
	"image/draw"
	"image/png"
	"io"
	"mime/multipart"
)

func ValidateSkin(skinFile io.Reader) error {
	skinFileCopy := skinFile
	img, err := png.Decode(skinFileCopy)
	if err != nil {
		return err
	}
	if img.Bounds().Dx() == 64 && (img.Bounds().Dy() == 32 || img.Bounds().Dy() == 64) {
		return nil
	}

	return fmt.Errorf("invalid skin dimensions (received=%dx%d, expected=64x32 or 64x64)", img.Bounds().Dx(), img.Bounds().Dy())
}

func RenderHeadOutOfSkin(skinFileHeader multipart.FileHeader) ([]byte, error) {
	buffer := new(bytes.Buffer)

	skinFile, err := skinFileHeader.Open()
	if err != nil {
		return buffer.Bytes(), err
	}
	defer skinFile.Close()

	fmt.Println("11")

	//if err := ValidateSkin(skinFile); err != nil {
	//	return buffer.Bytes(), err
	//}

	fmt.Println("22")

	skinImage, err := png.Decode(skinFile)
	if err != nil {
		return buffer.Bytes(), err
	}

	fmt.Println("33")

	siBounds := skinImage.Bounds()
	siNRGBA := image.NewNRGBA(image.Rect(0, 0, siBounds.Dx(), siBounds.Dy()))
	draw.Draw(siNRGBA, siNRGBA.Bounds(), skinImage, skinImage.Bounds().Min, draw.Src)

	fmt.Println("44")

	head := skin.RenderHead(siNRGBA, skin.Options{Scale: 20, Overlay: false, Square: true})

	err = png.Encode(buffer, head)

	fmt.Println("55")

	return buffer.Bytes(), err
}
