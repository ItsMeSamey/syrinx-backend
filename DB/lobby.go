package DB

import (
  "errors"
  
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
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

/// Get a lobby in which user is meant to be
/// Creates a new one it they are joining new
func LobbyFromUserSessionID(SessionID SessID) (*Lobby, error) {
  query := bson.M{
    "players": bson.M{
      "$elemMatch": bson.M{ "sessionID": SessionID, },
    },
  }

  result := LobbyDB.Coll.FindOne(LobbyDB.Context, query)
  err := result.Err()

  if err == mongo.ErrNoDocuments{
    // User are in no lobby, create a new one
    return createLobby(SessionID)
  } else if err != nil {
    return nil, errors.New("LobbyFromUserSessionID: DB.get error\n" + err.Error())
  }

  var lobby Lobby
  if err := result.Decode(&lobby); err != nil {
    return nil, errors.New("LobbyFromUserSessionID: Decode error\n" + err.Error())
  }
  return &lobby, nil
}

/// Adds a user and their team to a lobby if one exists or create a new one for them
func createLobby(SessionID SessID) (*Lobby, error) {
  /// Get the user object (to get the team ID)
  var user User
  if err := UserDB.get("sessionID", SessionID, &user); err != nil {
    return nil, errors.New("createLobby: UserDB.get error\n" + err.Error())
  }

  /// get user's team
  var team Team
  if err := TeamDB.get("teamID", user.TeamID, &team); err != nil {
    return nil, errors.New("createLobby: TeamDB.get error\n" + err.Error())
  }

  /// Get all the players in that team
  cursor, err := UserDB.Coll.Find(UserDB.Context, bson.M{"teamID": user.TeamID})
  if err != nil {
    return nil, errors.New("createLobby: Find error\n" + err.Error())
  }

  var players []Player
  if err = cursor.All(UserDB.Context, &players); err != nil {
    return nil, errors.New("createLobby: cursor.All error\n" + err.Error())
  }

  /// Try to find a partially vacant lobby
  result := LobbyDB.Coll.FindOne(LobbyDB.Context, bson.M{"isComplete": false})
  err = result.Err()

  /// Create a brand new lobby
  if err == mongo.ErrNoDocuments {
    lobby := Lobby{
      ID: nil,
      Players: players,
      Teams: []Team{team},
      IsComplete: false,
    }

    insertOneResult, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
    if err != nil {
      return nil, errors.New("createLobby: InsertOne error\n" + err.Error())
    }

    ID, ok := insertOneResult.InsertedID.(primitive.ObjectID)
    if !ok {
      return nil, errors.New("createLobby: Id was not object ID")
    }

    lobby.ID = &ID

    return &lobby, nil
  } else if err != nil {
    return nil, errors.New("createLobby: FindOne error\n" + err.Error())
  }

  /// Add to the existing lobby
  var lobby Lobby
  err = result.Decode(&lobby)
  if err != nil {
    return nil, errors.New("createLobby: result.Decode error\n" + err.Error())
  }

  lobby.Players = append(lobby.Players, players...)
  lobby.Teams = append(lobby.Teams, team)
  if len(lobby.Teams) >= 4 {
    lobby.IsComplete = true
  }

  _, err = LobbyDB.Coll.ReplaceOne(LobbyDB.Context, bson.M{"_id": lobby.ID}, lobby)
  if err != nil {
    return nil, errors.New("\n" + err.Error())
  }

  return &lobby, nil
}

func SaveLobby(lobby *Lobby) error {
  _, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
  return err
}

