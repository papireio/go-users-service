package database

type Session struct {
	Token     string `bson:"token"`
	CreatedAt string `bson:"createdAt"`
}
