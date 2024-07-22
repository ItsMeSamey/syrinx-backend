package DB

import (
  "errors"
  "crypto/rand"

  "go.mongodb.org/mongo-driver/bson"
)

/// User struct to store user information
type User struct {
  ID            ObjID  `bson:"_id,omitempty"`
  Username      string `bson:"user"`
  Email         string `bson:"mail"`
  Password      string `bson:"pass"`
  TeamID        TID    `bson:"teamID"`
  DiscordID     string `bson:"discordID"`
  SessionID     SessID `bson:"sessionID"`
  EmailReceived bool   `bson:"mailReceived"`
}

/// When we are creating users, we need a different struct than the above
type CreatableUser struct {
  Username  string  `bson:"user"`
  Email     string  `bson:"mail"`
  Password  string  `bson:"pass"`
  TeamID    TID     `bson:"teamID"`
  TeamName  *string `bson:"teamName"`
  DiscordID string  `bson:"discordID"`
}

/// Generates a unique SessionID
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
  exists, err := UserDB.exists(bson.M{"sessionID": bytes})
  if exists {
    if times > 1024 {
      return nil, errors.New("genSessionID: Lucky Error")
    }
    times += 1
    goto start
  }
  return ID, err
}

/// Generates a unique TeamID
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
  exists, err := UserDB.exists(bson.M{"teamID": bytes})
  if exists {
    if times > 1024*1024 {
      return nil, errors.New("genTeamID: OOPS! Good Luck!")
    }
    times += 1
    goto start
  }
  return ID, err
}

/// Function to create a user in DB
func CreateUser(user *CreatableUser) (SessID, error) {
  exists, err := UserDB.exists(bson.M{"user": user.Username})
  if err != nil { return nil, errors.New("CreateUser: Error while username lookup\n"+ err.Error()) }
  if exists { return nil, errors.New("CreateUser: User already exists") }

  exists, err = UserDB.exists(bson.M{"mail": user.Email})
  if err != nil { return nil, errors.New("CreateUser: Error while email lookup\n"+ err.Error()) }
  if exists { return nil, errors.New("CreateUser: Email cannot be reused") }

  exists, err = UserDB.exists(bson.M{"discordID": user.DiscordID})
  if err != nil { return nil, errors.New("CreateUser: Error while discordID lookup\n"+ err.Error()) }
  if exists { return nil, errors.New("CreateUser: Discord ID cannot be reused") }

  if user.TeamID == nil {
    if user.TeamName == nil || *(user.TeamName) == "" {
      return nil, errors.New("CreateUser: Team name needs to be specified")
    }
    tid, err := genTeamID()
    if err != nil {
      return nil, errors.New("CreateUser: Could not generate teamID\n"+ err.Error())
    }
    user.TeamID = tid

    err = createNewTeam(user)
    if err != nil {
      return nil, errors.New("CreateUser: Could not create team in db\n"+ err.Error())
    }
  } else {
    num, err := UserDB.Coll.CountDocuments(UserDB.Context, bson.M{"teamID": user.TeamID})
    if err != nil {
      return nil, errors.New("CreateUser: Error while team lookup\n"+ err.Error())
    }
    if num == 0 {
      return nil, errors.New("CreateUser: Team does not exist")
    } else if num >=4 {
      return nil, errors.New("CreateUser: Team already at max capacity")
    }

    name, err := TeamNameByID(user.TeamID)
    if err != nil {
      return nil, errors.New("CreateUser: could not get name of the team\n"+ err.Error())
    }
    user.TeamName = &name;
  }

  SessionID, err := genSessionID()
  if err != nil { return nil, err }

  go sendEmailAsync(user)

  _, err = UserDB.Coll.InsertOne(UserDB.Context, &User{
    Username     : user.Username,
    Email        : user.Email,
    Password     : user.Password,
    TeamID       : user.TeamID,
    DiscordID    : user.DiscordID,
    SessionID    : SessionID,
    EmailReceived: false,
  })
  return SessionID, err
}

/// The user authantication function
func UserAuthenticate(username, password string) (*User, error) {
  var user User
  result := UserDB.Coll.FindOne(UserDB.Context, bson.M{"user": username, "pass":password})
  if result == nil {
    return nil, errors.New("UserAuthenticate: Invalid Password/Username")
  }
  err := result.Decode(&user)
  if err != nil {
    return nil, err
  }
  return &user, err
}

/// Get Uset object from session id
func UserFromSessionID(SessionID SessID) (*User, error) {
  var user User
  return &user, UserDB.get(bson.M{"sessionID": SessionID}, &user)
}

