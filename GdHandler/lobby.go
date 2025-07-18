package GdHandler

import (
  "errors"
  "net/http"
  "strconv"
  "sync"
  
  "ccs.ctf/DB"
  
  "github.com/gorilla/websocket"
  "go.mongodb.org/mongo-driver/bson"
)

type (
  Player struct {
    ID        DB.ObjID        `bson:"_id,omitempty"`
    Username  string          `bson:"user"`
    Email     string          `bson:"mail"`
    DiscordID string          `bson:"discordID"`
    SessionID DB.SessID       `bson:"sessionID"`
    Conn      *websocket.Conn `bson:"-"`
  }

  Lobby struct {
    Team        *DB.Team
    Players     []*Player
    PlayerMutex sync.RWMutex
    Upgrader    websocket.Upgrader
  }
)

func makeLobby(ID DB.TID) (*Lobby, error) {
  team, err := DB.TeamByTeamID(ID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.TeamByID error\n" + err.Error())
  }

  if !team.Exception && team.Level != LEVEL {
    return nil, errors.New("LobbyIDFromUserSessionID: Team Level Mismatch")
  }

  lobby := makeLobbyFromTeam(team)

  err = lobby.populatePlayers()
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: lobby.populatePlayers error\n" + err.Error())
  }

  return lobby, nil
}

func makeLobbyFromTeam(team *DB.Team) *Lobby {
  return &Lobby {
    Team: team,
    PlayerMutex: sync.RWMutex{},
    Upgrader: websocket.Upgrader{
      ReadBufferSize:  1024,
      WriteBufferSize: 1024,
      CheckOrigin:     func (r *http.Request) bool {return true},
    },
  }
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

  lobby.Players = players

  return nil
}

func LobbyIDFromUserSessionID(SessionID DB.SessID) (DB.TID, int, error) {
  user, err := DB.UserFromSessionID(SessionID)
  if err != nil {
    return nil, 0, errors.New("LobbyIDFromUserSessionID: DB.UserFromSessionID error\n" + err.Error())
  }

  lobby, err := getAddedLobby(user.TeamID)
  if err != nil {
    return nil, 0, errors.New("LobbyIDFromUserSessionID: getAddedLobby error\n" + err.Error())
  }
  if lobby.Team.Level == LEVEL || lobby.Team.Exception {
    return lobby.Team.TeamID, lobby.Team.Level, nil
  }
  return nil, 0, errors.New("LobbyIDFromUserSessionID error: player of level " + strconv.Itoa(lobby.Team.Level) + " cannot join level " + strconv.Itoa(LEVEL))
}

/// Forcefully close the lobby
func (lobby *Lobby) delete() {
  lobby.PlayerMutex.Lock()
  defer lobby.PlayerMutex.Unlock()

  for i := range lobby.Players {
    if lobby.Players[i].Conn != nil {
      lobby.Players[i].Conn.Close()
      lobby.Players[i].Conn = nil
    }
  }
}

