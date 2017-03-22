package data

import "github.com/dimorinny/twitch-interesting-fragments/configuration"

const (
	mongoStorage = "mongo"
)

type Storage interface {
	AddUploadedFragment(fragment UploadedFragment) error
}

func InitStorage(configuration configuration.Configuration) (Storage, error) {
	switch configuration.StorageType {
	case mongoStorage:
		return NewMongoStorage(configuration.StorageHost)
	default:
		return NewStubStorage(), nil
	}
}
