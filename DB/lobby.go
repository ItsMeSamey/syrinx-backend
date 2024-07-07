package DB

import (
  "errors"
)

/// Struct meant to be used in GdIntegration
type Player struct {
  ID        ObjID       `bson:"_id,omitempty"`
  SessionID SessID      `bson:"sessionID"`
  IN        chan []byte `bson:"-"`
}

type Lobby struct {
  ID    ObjID    `bson:"_id,omitempty"`
  Users []Player `bson:"users"`
}

func GetLobby(lobbyID ObjID) ([]Player, error) {
  var lobby Lobby
  err := LobbyDB.get("_id", lobbyID, &lobby)
  if err != nil {
    return nil, err
  }
  if len(lobby.Users) == 0 {
    return nil, errors.New("GetLobby: Lobby is empty")
  }
  return lobby.Users, nil
}

func SaveLobby(lobby *Lobby) error {
  _, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
  return err
}

