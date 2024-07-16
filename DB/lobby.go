package DB

import (
  "errors"
  
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
)

/// Struct meant to be used in GdIntegration
type Player struct {
  ID        ObjID       `bson:"_id,omitempty"`
  SessionID SessID      `bson:"sessionID"`
  TeamID    TID         `bson:"teamID"`
  IN        chan []byte `bson:"-"`
}

type Lobby struct {
  ID          ObjID   `bson:"_id,omitempty"`
  Players    []Player `bson:"players"`
  Teams      []Team   `bson:"teams"`
  IsComplete bool     `bson:"isComplete"`
}

func LobbyFromID(lobbyID ObjID) (*Lobby, error) {
  var lobby Lobby
  
  if err := LobbyDB.get("_id", lobbyID, &lobby); err != nil {
    return nil, errors.New("LobbyFromID: DB.get error\n" + err.Error())
  }

  return &lobby, nil
}

func LobbyFromUserSessionID(SessionID SessID) (*Lobby, error) {
  var lobby *Lobby
  query := bson.M{
    "users": bson.M{
      "$elemMatch": bson.M{ "sessionID": SessionID, },
    },
  }

  result := LobbyDB.Coll.FindOne(LobbyDB.Context, query)
  err := result.Err()

  if err == mongo.ErrNoDocuments{
    lobby, err = createLobby(SessionID)
  } else if err != nil {
    return nil, errors.New("LobbyFromUserSessionID: DB.get error\n" + err.Error())
  } else{
    var _lobby Lobby
    if err := result.Decode(&_lobby); err != nil {
      return nil, errors.New("LobbyFromUserSessionID: Decode error\n" + err.Error())
    }
    lobby = &_lobby
  }
  return lobby, nil
}

func createLobby(SessionID SessID) (*Lobby, error) {
  var user User
  if err := UserDB.get("SessionID", SessionID, &user); err != nil {
    return nil, errors.New("createLobby: UserDB.get error\n" + err.Error())
  }

  var team Team
  if err := TeamDB.get("teamID", user.TeamID, &team); err != nil {
    return nil, errors.New("createLobby: TeamDB.get error\n" + err.Error())
  }
  
  return nil, errors.ErrUnsupported
}

func SaveLobby(lobby *Lobby) error {
  _, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
  return err
}

