package data

type UploadedFragment struct {
	ID   string  `bson:"_id"`
	Url  string  `bson:"url"`
	Name string  `bson:"name"`
	Time int64   `bson:"time"`
	Rate float32 `bson:"rate"`
}
