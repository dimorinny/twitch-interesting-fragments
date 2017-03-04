package api

import (
	"fmt"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	"net/http"
	"encoding/json"
)

const (
	scheme = "http"
)

type (
	Uploader struct {
		configuration configuration.Configuration
		client        *http.Client
	}
)

func NewUploader(configuration configuration.Configuration, client *http.Client) *Uploader {
	return &Uploader{
		configuration: configuration,
		client:        client,
	}
}

func (u *Uploader) Upload() (*UploadResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s://%s:%d/api/v1/upload",
			scheme,
			u.configuration.UploaderHost,
			u.configuration.UploaderPort,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &UploadResponse{}
	json.NewDecoder(resp.Body).Decode(response)

	return response, nil
}
