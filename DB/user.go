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
  DiscordID string `bson:"discordID"`
  SessionID SessID `bson:"sessionID"`
  EmailReceived bool `bson:"mailReceived"`
}

type CreatableUser struct {
  Username  string `bson:"user"`
  Email   string `bson:"mail"`
  Password  string `bson:"pass"`
  TeamID  TID  `bson:"teamID"`
  TeamName *string `bson:"teamName"`
  DiscordID string `bson:"discordID"`
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
func internalSendConfirmationEmail(user *CreatableUser) error {
  const subject = "Subject: Confirmation for participation in Syrinx\n"
  const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

  tmpl, err := template.ParseFiles("email_template.html")
  if err != nil {
    return err
  }

  var body bytes.Buffer
  err = tmpl.Execute(&body, struct {
    JoinedMember string
    Email        string
    TeamName     string
    TeamID       string
  }{
    JoinedMember: user.Username,
    Email:        user.Email,
    TeamName:     *user.TeamName,
    TeamID:       hex.EncodeToString(user.TeamID[:]),
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

func internalUpdateEmailStatus(user *CreatableUser) error {
  _, err := UserDB.Coll.UpdateOne(UserDB.Context, bson.D{{"user", user.Username}}, bson.D{{"$set", bson.D{{"mailReceived", true}}}})
  return err
}

/// Must run this as Async
func sendEmailAsync(user *CreatableUser) {
  const maxCount = 60*12
  var err error
  count := 0

  if user.TeamName == nil || *user.TeamName == "" {
    var team Team
    err = TeamDB.get("teamID", user.TeamID, &team)
    for err != nil && count < maxCount {
      count += 1
    }
    user.TeamName = &(team.TeamName);
  }

  err = internalSendConfirmationEmail(user)
  for err != nil && count < maxCount {
    count += 1
    time.Sleep(time.Minute)
    err = internalSendConfirmationEmail(user)
  }

  err = internalUpdateEmailStatus(user)
  for err != nil && count < maxCount {
    count += 1
    time.Sleep(time.Minute)
    err = internalUpdateEmailStatus(user)
  }
}

func CreateUser(user *CreatableUser) (SessID, error) {

  exists, err := UserDB.exists("user", user.Username)
  if err != nil { return nil, err }
  if exists { return nil, errors.New("CreateUser: User already exists") }

  exists, err = UserDB.exists("mail", user.Email)
  if err != nil { return nil, err }
  if exists { return nil, errors.New("CreateUser: Email cannot be reused") }

  exists, err = UserDB.exists("discordID", user.DiscordID)
  if err != nil { return nil, err }
  if exists { return nil, errors.New("CreateUser: Discord ID cannot be reused") }

  if user.TeamID == nil {
    tid, err := genTeamID()
    if err != nil {
      return nil, err
    }
    user.TeamID = tid
  } else {
    num, err := UserDB.Coll.CountDocuments(UserDB.Context, bson.D{{"teamID", user.TeamID}})
    if err != nil {
      return nil, err
    }
    if num == 0 {
      return nil, errors.New("CreateUser: Team does not exist")
    } else if num >=4 {
      return nil, errors.New("CreateUser: Team already at max capacity")
    }
  }

  SessionID, err := genSessionID()
  if err != nil { return nil, err }

  go sendEmailAsync(user)

  _, err = UserDB.Coll.InsertOne(UserDB.Context, *user)
  return SessionID, err
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

