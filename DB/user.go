package DB

import (
	"crypto/rand"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// User struct to store user information
type User struct {
	ID        string `bson:"_id,omitempty"`
	Username  string `bson:"user"`
	Email     string `bson:"mail"`
	Password  string `bson:"pass"`
	TeamID    string `bson:"teamID"`
	DiscordID string `bson:"discordID"`
	SessionID SessID `bson:"SessionID"`
}

func genSessionID() (SessID, error) {
	bytes := make([]byte, 6*64)
	_, err := rand.Read(bytes)
	ID := SessID(bytes)
	if ID == nil {
		return nil, errors.New("genSessionID: ID generation failed")
	}
	return ID, err
}

func CreateUser(user *User) (SessID, error) {
	if user.ID != "" {
		return nil, errors.New("CreateUser: ID cannot be set")
	}
	if user.SessionID != nil {
		return nil, errors.New("CreateUser: SessionID cannot be set")
	}

	ID, err := genSessionID()
	if err != nil {
		return nil, err
	}

	user.SessionID = ID

	exists, err := UserDB.exists("user", user.Username)
	if exists {
		return nil, errors.New("CreateUser: User exists")
	}
	if err != nil {
		return nil, err
	}

	_, err = UserDB.coll.InsertOne(UserDB.context, *user)
	if err != nil {
		return nil, err
	}

	return ID, nil
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

func UserFromSessionID(SessionID SessID) (*User, error) {
	var user User
	return &user, UserDB.get("SessionID", SessionID, &user)
}

