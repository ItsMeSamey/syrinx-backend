package DB

import (
  "errors"
  "log"
  
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
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
  Teams      []Team   `bson:"-"`
  IsComplete bool     `bson:"isComplete"`
}

const MAX_TEAMS = 4


func (lobby *Lobby) PopulateTeams() error {
  teams := []Team{}
  for _, player := range lobby.Players {
    exists := false
    for _, _team := range teams {
      if *(_team.TeamID) == *(player.TeamID) {
        exists = true
        break
      }
    }
    if exists {
      continue
    }

    var team Team
    err := TeamDB.get(bson.M{"teamID": player.TeamID}, &team)
    if err != nil {
      return errors.New("populateTeams: DB.get error\n" + err.Error())
    }
    teams = append(teams, team)
  }
  return nil
}

/// Convert LobbyTemplate to a real lobby
func (lobby *Lobby) CreateInsert() error {
  if lobby.ID != nil {
    return errors.New("Lobby.CreateInsert: called on already existing lobby")
  }

  insertOneResult, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
  if err != nil {
    return errors.New("Lobby.CreateInsert: InsertOne error\n" + err.Error())
  }

  ID, ok := insertOneResult.InsertedID.(primitive.ObjectID)
  if !ok {
    return errors.New("Lobby.CreateInsert: Id was not object ID")
  }

  lobby.ID = &ID
  lobby.IsComplete = len(lobby.Teams) >= MAX_TEAMS || lobby.IsComplete

  return nil
}

func (lobby *Lobby) Sync(maxTries byte) error {
  return LobbyDB.syncTryHard(bson.M{"_id": lobby.ID}, lobby, maxTries)
}

func (lobby *Lobby) Merge(lobbyTemplate *Lobby) {

  if lobby.IsComplete {
    log.Println("Merge called on complete lobby !!")
  }

  /// Add to the existing lobby
  lobby.Players = append(lobby.Players, lobbyTemplate.Players...)
  lobby.Teams = append(lobby.Teams, lobbyTemplate.Teams...)

  lobby.IsComplete = len(lobby.Teams) >= MAX_TEAMS || lobby.IsComplete
}

func LobbyFromID(lobbyID ObjID) (*Lobby, error) {
  var lobby Lobby
  
  if err := LobbyDB.get(bson.M{"_id": lobbyID}, &lobby); err != nil {
    return nil, errors.New("LobbyFromID: DB.get error\n" + err.Error())
  }

  return &lobby, nil
}

func GetIncompleteLobby() (*Lobby, bool, error) {
  var lobby Lobby
  
  exists, err := LobbyDB.getExists(bson.M{"isComplete": false}, &lobby)

  if err != nil {
    return nil, false, errors.New("Lobby.Sync error\n" + err.Error())
  }
  
  if !exists {
    return nil, false, nil
  }

  return &lobby, true, nil
}

/// Get a lobby in which user is meant to be
func LobbyFromUserSessionID(SessionID SessID) (*Lobby, bool, error) {
  var lobby Lobby

  exists, err := LobbyDB.getExists(bson.M{
    "players": bson.M{
      "$elemMatch": bson.M{ "sessionID": SessionID, },
    },
  }, &lobby)

  if err != nil {
    return nil, false, errors.New("LobbyFromUserSessionID error\n" + err.Error())
  }
  
  if !exists {
    return nil, false, nil
  }

  return &lobby, true, nil
}

/// Create a lobby template with nil ID
func NewLobbyTemplate(SessionID SessID) (*Lobby, error) {
  /// Get the user object (to get the team ID)
  var user User
  if err := UserDB.get(bson.M{"sessionID": SessionID}, &user); err != nil {
    return nil, errors.New("NewLobbyTemplate: UserDB.get error\n" + err.Error())
  }

  /// get user's team
  var team Team
  if err := TeamDB.get(bson.M{"teamID": user.TeamID}, &team); err != nil {
    return nil, errors.New("NewLobbyTemplate: TeamDB.get error\n" + err.Error())
  }

  /// Get all the players in that team
  cursor, err := UserDB.Coll.Find(UserDB.Context, bson.M{"teamID": user.TeamID})
  if err != nil {
    return nil, errors.New("NewLobbyTemplate: Find error\n" + err.Error())
  }

  var players []Player
  if err = cursor.All(UserDB.Context, &players); err != nil {
    return nil, errors.New("NewLobbyTemplate: cursor.All error\n" + err.Error())
  }

  /// Try to find a partially vacant lobby
  return &Lobby{
    ID: nil,
    Players: players,
    Teams: []Team{team},
    IsComplete: 1 == MAX_TEAMS,
  }, nil
}

