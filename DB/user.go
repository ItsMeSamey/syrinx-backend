package DB

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// User struct to store user information
type User struct {
	Username  string `bson:"user"`
	Password  string `bson:"pass"`
	TeamID    int    `bson:"teamID"`
	DiscordID string `bson:"discordID"`
}

// func genSessionID() (string, error) {
// 	bytes := make([]byte, 6*16)
// 	if _, err := rand.Read(bytes); 
// 	err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(bytes), nil
// }

func CreateUser(user *User) error {
	result := UserDB.coll.FindOne(UserDB.context, bson.D{{"user", user.Username}})
	if result == nil {
		return errors.New("CreateUser: Result is `nil`")
	}
	err := result.Err()
	if err == nil {
		return errors.New("CreateUser: User exists")
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	insert, err := UserDB.coll.InsertOne(UserDB.context, *user)
	_ = insert
	
	if err != nil {
		return err
	}
	return nil
}

func UserAuthenticate(username, password string) (*User, error) {
	var user User
	result := UserDB.coll.FindOne(UserDB.context, bson.D{{"user", username}, {"pass", password}})
	if result == nil {
		return nil, errors.New("UserAuthenticate: Invalid Password")
	}
	err := result.Decode(user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func UserFromSessionID(_id string) (*User, error) {
	var user User
	result := UserDB.coll.FindOne(UserDB.context, bson.D{{"_id", _id}})
	if result == nil {
		return nil, errors.New("UserFromSessionID: Token")
	}
	err := result.Decode(user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

