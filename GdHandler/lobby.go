package GdHandler

import (
  "errors"
  "sync"
  
  "ccs.ctf/DB"
  
  "github.com/gorilla/websocket"
  "go.mongodb.org/mongo-driver/bson"
)

type Player struct {
  ID        DB.ObjID    `bson:"_id,omitempty"`
  SessionID DB.SessID   `bson:"sessionID"`
  IN        chan []byte `bson:"-"`
}

type Lobby struct {
  Team        *DB.Team
  Players     []*Player
  Playercount byte
  PlayerMutex sync.RWMutex
  Upgrader    websocket.Upgrader
  Deadtime    byte
}

func makeLobby(ID DB.TID) (*Lobby, error) {
  team, err := DB.TeamByTeamID(ID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.TeamByID error\n" + err.Error())
  }

  lobby := &Lobby {
    Team: team,
    Playercount: 0,
    PlayerMutex: sync.RWMutex{},
    Upgrader: websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024, CheckOrigin: originChecker},
    Deadtime: 0,
  }

  lobby.populatePlayers()
  return lobby, nil
}

func (lobby *Lobby) populatePlayers() error {
  /// Get all the players in that team
  cursor, err := DB.UserDB.Coll.Find(DB.UserDB.Context, bson.M{"teamID": lobby.Team.TeamID})
  if err != nil {
    return errors.New("NewLobbyTemplate: Find error\n" + err.Error())
  }

  var players []*Player
  if err = cursor.All(DB.UserDB.Context, &players); err != nil {
    return errors.New("NewLobbyTemplate: cursor.All error\n" + err.Error())
  }

  return nil
}

func LobbyIDFromUserSessionID(SessionID DB.SessID) (DB.TID, error) {
  user, err := DB.UserFromSessionID(SessionID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.UserFromSessionID error\n" + err.Error())
  }

  team, err := DB.TeamByTeamID(user.TeamID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.TeamByID error\n" + err.Error())
  }

  return team.TeamID, nil
}

