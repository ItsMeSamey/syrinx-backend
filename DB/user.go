package DB

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// User struct to store user information
type User struct {
	Username  string `bson:"user"`
	Email     string `bson:"mail"`
	Password  string `bson:"pass"`
	TeamID    string `bson:"teamID"`
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
	exists, err := UserDB.exists("user", user.Username)
	if exists {
		return errors.New("CreateUser: User exists")
	}
	if err != nil {
		return err
	}

	result, err := UserDB.coll.InsertOne(UserDB.context, *user)
	_ = result
	
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
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func UserFromSessionID(_id string) (*User, error) {
	var user User
	return &user, UserDB.get("_id", _id, &user)
}

