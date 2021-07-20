package internal

import (
	"encoding/base64"
	"encoding/json"
)

type FakeDB map[string]Images
type Images map[int]*Image

type authToken struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

type Image struct {
	Mime string `json:"mime"`
	Data string `json:"data"`
}

// NewImageJson Make sure that data in Image are valid.
func NewImageJson(data *[]byte) (image *Image, err error) {
	err = json.Unmarshal(*data, &image)
	if err != nil {
		return
	}
	_, err = base64.URLEncoding.DecodeString(image.Data)
	if err != nil {
		return
	}
	return
}
