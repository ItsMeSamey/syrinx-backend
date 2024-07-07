package DB

import (
	"crypto/rand"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// User struct to store user information
type User struct {
	ID        ObjID  `bson:"_id,omitempty"`
	Username  string `bson:"user"`
	Email     string `bson:"mail"`
	Password  string `bson:"pass"`
	TeamID    TID    `bson:"teamID"`
	DiscordID string `bson:"discordID"`
	SessionID SessID `bson:"SessionID"`
}

func genSessionID() (SessID, error) {
	times := 0
	start:

	bytes := make([]byte, 6*64)
	_, err := rand.Read(bytes)
	ID := SessID(bytes)
	if ID == nil {
		return nil, errors.New("genSessionID: ID generation failed")
	}
	if err != nil {
		return nil, err
	}
	exists, err := UserDB.exists("SessionID", bytes)
	if exists {
		if times > 1024 {
			return nil, errors.New("genSessionID: I'm a Bathtub, (if you are seeing this, contact us IMMEDIATELY!)")
		}
		times += 1
		goto start
	}
	return ID, err
}

func genTeamID() (TID, error) {
	times := 0
	start:

	bytes := make([]byte, 3)
	_, err := rand.Read(bytes)
	ID := TID(bytes)
	if ID == nil {
		return nil, errors.New("genTeamID: ID generation failed")
	}
	if err != nil {
		return nil, err
	}
	exists, err := UserDB.exists("TeamID", bytes)
	if exists {
		if times > 1024*1024 {
			return nil, errors.New("genTeamID: OOPS, Lucky Draw!!, (if you are seeing this, contact us!)")
		}
		times += 1
		goto start
	}
	return ID, err
}

func CreateUser(user *User) error {
	if user.ID != nil {
		return errors.New("CreateUser: ID cannot be set")
	}
	if user.SessionID != nil {
		return errors.New("CreateUser: SessionID cannot be set")
	}

	SessionID, err := genSessionID()
	if err != nil {
		return err
	}

	user.SessionID = SessionID

	if user.TeamID == nil {
		tid, err := genTeamID()
		if err != nil {
			return err
		}
		user.TeamID = tid
	} else {
		exists, err := UserDB.exists("TeamID", user.TeamID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("CreateUser: Team does not exist")
		}
	}

	exists, err := UserDB.exists("user", user.Username)
	if exists {
		return errors.New("CreateUser: User exists")
	}
	if err != nil {
		return err
	}

	_, err = UserDB.coll.InsertOne(UserDB.context, *user)

	return err
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

