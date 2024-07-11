package DB

import (
  "crypto/rand"
  "errors"
  "time"
  
  "bytes"
  "encoding/hex"
  "html/template"
  "net/smtp"
  
  "go.mongodb.org/mongo-driver/bson"
)

// User struct to store user information
type User struct {
  ID    ObjID  `bson:"_id,omitempty"`
  Username  string `bson:"user"`
  Email   string `bson:"mail"`
  Password  string `bson:"pass"`
  TeamID  TID  `bson:"teamID"`
  TeamName  string  `bson:"teamName"`
  DiscordID string `bson:"discordID"`
  SessionID SessID `bson:"sessionID"`
  EmailReceived bool `bson:"mailReceived"`
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

/// blocking send email function
func internalSendConfirmationEmail(user *User) error {
  const subject = "Subject: Confirmation for participation in Syrinx\n"
  const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

  tmpl, err := template.ParseFiles("email_template.html")
  if err != nil {
    return err
  }

  var body bytes.Buffer
  err = tmpl.Execute(&body, struct {
    Username  string
    Email   string
    TeamID string
  }{
    Username:  user.Username,
    Email:   user.Email,
    TeamID: hex.EncodeToString(user.TeamID[:]),
  })
  if err != nil {
    return err
  }

  message := subject + mime + body.String()

  err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", EMAIL_SENDER, EMAIL_SENDER_PASSWORD, "smtp.gmail.com"),
                      EMAIL_SENDER, []string{user.Email}, []byte(message),
  )
  if err != nil {
    return err
  }

  return nil
}

func internalUpdateEmailStatus(user *User) error {
  _, err := UserDB.Coll.UpdateOne(UserDB.Context, bson.D{{"user", user.Username}}, bson.D{{"$set", bson.D{{"mailReceived", true}}}})
  return err
}

/// Must run this as Async
func sendEmailAsync(user *User) {
  var err error

  err = internalSendConfirmationEmail(user)
  for err != nil {
    time.Sleep(time.Minute)
    err = internalSendConfirmationEmail(user)
  }

  err = internalUpdateEmailStatus(user)
  for err != nil {
    time.Sleep(time.Minute)
    err = internalUpdateEmailStatus(user)
  }
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

  go sendEmailAsync(user)
  if err != nil {
    return err
  }

  _, err = UserDB.Coll.InsertOne(UserDB.Context, *user)
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

