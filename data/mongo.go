package data

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	db         = "twitch"
	collection = "fragments"
)

type MongoStorage struct {
	session *mgo.Session
}

func NewMongoStorage(session *mgo.Session) Storage {
	return &MongoStorage{
		session: session,
	}
}

func (s *MongoStorage) AddUploadedFragment(fragment UploadedFragment) error {
	recognitions := s.session.DB(db).C(collection)

	fragment.ID = bson.NewObjectId().Hex()
	fragment.Time = time.Now().Unix()

	return recognitions.Insert(fragment)
}
