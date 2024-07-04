package DB

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	userBucket = "users"
	teamBucket = "teams"
	sessionBucket = "sessions"
)

// User struct to store user information
type User struct {
	Username  string `json:"user"`
	Password  string `json:"pass"`
	UserId    string `json:"userID"`
	TeamID    int    `json:"teamID"`
	SessionID string `json:"sesisonID"`
}

func genSessionID() (string, error) {
	bytes := make([]byte, 6*16)
	if _, err := rand.Read(bytes); 
	err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (user *User) Create() error {
	// Return error if user is already present
	exists, err := UserDB.DoesExist(userBucket, []byte(user.Username))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CreateUser: User Exists")
	}

	tries := 0
start:
	sessionID, err := genSessionID()
	if err != nil {
		return err
	}

	exists, err = UserDB.DoesExist(sessionBucket, []byte(sessionID))
	if err != nil {
		return err
	}
	if exists {
		tries += 1
		if tries > 1024*1024 {
			return errors.New("Create: I'm a Teapot")
		}
		goto start
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return UserDB.addToBucket(userBucket, []byte(user.Username), data)
}

func (user *User) Authenticate() error {
	var tempUser User
	val, err := UserDB.getFromBucket(userBucket, []byte(user.Username))
	if err != nil {
		return err
	}

	err = json.Unmarshal(val, &tempUser)
	if err != nil {
		return err
	}

	if tempUser.Password != user.Password {
		return fmt.Errorf("Authenticate: Invalid Password")
	}
	user = &tempUser
	return nil
}

func (user *User) GetUserFromSessionID() error {
	if user.SessionID == "" {
		return errors.New("GetUserFromSessionID: SessionID not given")
	}

	val, err := UserDB.getFromBucket(sessionBucket, []byte(user.SessionID))
	if err != nil {
		return err
	}
	return json.Unmarshal(val, &user)
}


