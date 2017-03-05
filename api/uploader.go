package api

import (
	"encoding/json"
	"fmt"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	"github.com/kataras/go-errors"
	"net/http"
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
	response := &uploadResult{}

	err := u.do("upload", response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (u *Uploader) do(method string, response interface{}) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s://%s:%d/api/v1/%s",
			scheme,
			u.configuration.UploaderHost,
			u.configuration.UploaderPort,
			method,
		),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Response status code different from 200")
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
