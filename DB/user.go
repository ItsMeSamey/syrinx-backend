package DB

import (
	"crypto/rand"
	"errors"
	
    "encoding/hex"
	"net/smtp"
	"html/template"
	"bytes"
    "fmt"

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
	SessionID SessID `bson:"sessionID"`
}

func genSessionID() (SessID, error) {
	times := 0
	start:

	bytes := make([]byte, 6*64)
	_, err := rand.Read(bytes)
	/// This pacnics needlessly ??
	// if len(bytes) != 64 {
	// 	return nil, errors.New("genSessionID: a length mismatch happened. Panic avoided!!")
	// }
	ID := SessID(bytes)
	if ID == nil {
		return nil, errors.New("genSessionID: ID generation failed")
	}
	if err != nil {
		return nil, err
	}
	exists, err := UserDB.exists("sessionID", bytes)
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
	/// This pacnics needlessly
	// if len(bytes) != 3 {
	// 	return nil, errors.New("genTeamID: a length mismatch happened. Panic avoided!!")
	// }
	ID := TID(bytes)
	if ID == nil {
		return nil, errors.New("genTeamID: ID generation failed")
	}
	if err != nil {
		return nil, err
	}
	exists, err := UserDB.exists("teamID", bytes)
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
	user.ID = nil
	user.SessionID = nil

	exists, err := UserDB.exists("user", user.Username)
	if err != nil { return err }
	if exists { return errors.New("CreateUser: User already exists") }

	exists, err = UserDB.exists("mail", user.Email)
	if err != nil { return err }
	if exists { return errors.New("CreateUser: Email cannot be reused") }

	exists, err = UserDB.exists("discordID", user.DiscordID)
	if err != nil { return err }
	if exists { return errors.New("CreateUser: Discord ID cannot be reused") }

	if user.TeamID == nil {
		tid, err := genTeamID()
		if err != nil {
			return err
		}
		user.TeamID = tid
	} else {
		exists, err := UserDB.exists("teamID", user.TeamID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("CreateUser: Team does not exist")
		}
	}

	SessionID, err := genSessionID()
	if err != nil { return err }

	user.SessionID = SessionID

	_, err = UserDB.Coll.InsertOne(UserDB.Context, *user)
	sendConfirmationEmail(user)
	return err
}

func UserAuthenticate(username, password string) (*User, error) {
	var user User
	result := UserDB.Coll.FindOne(UserDB.Context, bson.D{{"user", username}, {"pass", password}})
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
	return &user, UserDB.get("sessionID", SessionID, &user)
}





func sendConfirmationEmail(user *User) error {
    tmpl, err := template.ParseFiles("email_template.html")
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    // Convert TeamID ([]int) to []byte
    teamIDBytes := make([]byte, len(user.TeamID))
    for i, v := range user.TeamID {
        teamIDBytes[i] = byte(v)
    }

    // Convert []byte to hex string
    hexStr := hex.EncodeToString(teamIDBytes)

    var body bytes.Buffer
    err = tmpl.Execute(&body, struct {
        Username  string
        Email     string
        TeamIDHex string
    }{
        Username:  user.Username,
        Email:     user.Email,
        TeamIDHex: hexStr,
    })
    if err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    from := "riteshkapoor1314@gmail.com"
    to := []string{user.Email}
    subject := "Subject: Confirmation for participation in Syrinx\n"
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    message := subject + mime + body.String()

    err = smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, "cpkmfnxfrjjqtysy", "smtp.gmail.com"),
        from, to, []byte(message))
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}