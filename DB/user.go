package DB

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	userBucket = "users"
	teamBucket = "teams"
	sessionsBucket = "sessions"
)

// User struct to store user information
type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	UserId string	  `json:"userID"`
	TeamID   int    `json:"teamID"`
	SessionID string `json:"sesisonID"`
}

func setSessionID(user *User) (error) {
	// Set the user's ssession id to a <unique> and random base64 encoded string
	// also make a bucket to hold session keys and respective user names
	// TODO: when user reauthanticates, old one should be deleted and a new token must be generated
	_ = user
	return nil
}


//start*****

func generateSessionID() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	if _, err := rand.Read(bytes); 
	err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func createSession(db *bolt.DB, userID, username string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// Create or get the sessions bucket
		sessionsBucket, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
		if err != nil {
			return err
		}

		// Create or get the user bucket
		userBucket, err := tx.CreateBucketIfNotExists([]byte(userBucket))
		if err != nil {
			return err
		}

		// Store session ID associated with user ID
		err = sessionsBucket.Put([]byte(sessionID), []byte(userID))
		if err != nil {
			return err
		}

		// Store session ID associated with username
		err = usernamesBucket.Put([]byte(sessionID), []byte(username))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return sessionID, nil
}
func deleteExistingSession(tx *bolt.Tx, userID string) error {
	userToSessionBucket := tx.Bucket([]byte(userIDToSessionBucket))
	if userToSessionBucket == nil {
		// No existing sessions for this userID
		return nil
	}
	//getting existing session id
	get := func(tx *bolt.Tx) err{
		bucket = tx.Bucket([]byte(""))
		val := bucket.Get([]byte(""))
		if val == nil{
			//not found
			return err.New("No existing session")
		}
		fmt.Println(val)
		return nil
	}
	if err :=db.View(get);
	err != nill{
		log.Fatal(err)
	}

}

//end*******


func UserExists(username string) (bool, error) {
	// Implement this
	return true, nil
}

func CreateUser(user *User) error {
	// Return error if user is already present
	exists, err := UserExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CreateUser: User Exists")
	}
	return DBInstance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(user.Username), data); err != nil {
			return err
		}
		return nil
	})
}

func Authenticate(username, password string) (*User, error) {
	var user User
	err := DBInstance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		val := b.Get([]byte(username))
		if val == nil {
			return fmt.Errorf("Authenticate: User not found")
		}
		return json.Unmarshal(val, &user)
	})

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("Authenticate: Invalid Password")
	}
	
	err = setSessionID(&user)
	if err != nil {
		return nil, fmt.Errorf("Authenticate: SessionID Creation Failed")
	}
	return &user, nil
}

func GetUserFromSessionID(sessionID string) (*User, error) {
	// TODO: lookup the sesisonID table

	return nil, nil
}

//for eg
UserID= "1"

func GetSessionIDFromUser(user *User) (string, error) {
	// TODO: Implement creation of user's sessionID 
	
		sessionID,err := createSession(db, userID)
		if err !=nil{
			log.Fatal(err)
		}
		fmt.Printf("Created session ID: %s for userID: %s\n", &sessionID,user.UserId)

	// NOTE: One user must have only 1 session ID,

	return "", nil
}

