package database

type User struct {
	Name            string    `bson:"name"`
	Email           string    `bson:"email"`
	Uuid            string    `bson:"uuid"`
	PasswordHash    string    `bson:"hash"`
	ValidationToken string    `bson:"validation_token"`
	Sessions        []Session `bson:"sessions"`
}
