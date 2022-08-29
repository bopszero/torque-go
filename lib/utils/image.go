package utils

import (
	"bytes"
	"image"
	"image/jpeg"
)

func ConvertImageToBytes(imageDecoded image.Image) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, imageDecoded, nil)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
