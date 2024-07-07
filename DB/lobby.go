package DB

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

/// Struct meant to be used in GdIntegration
type Player struct {
  ID        ObjID       `bson:"_id,omitempty"`
  SessionID SessID      `bson:"sessionID"`
  IN        chan []byte `bson:"-"`
}

type Lobby struct {
  ID    ObjID    `bson:"_id,omitempty"`
  Players []Player `bson:"users"`
}

func LobbyFromID(lobbyID ObjID) (*Lobby, error) {
  var lobby Lobby
  err := LobbyDB.get("_id", lobbyID, &lobby)
  if err != nil {
    return nil, err
  }
  if len(lobby.Players) == 0 {
    return nil, errors.New("GetLobby: Lobby is empty")
  }
  return &lobby, nil
}

func LobbyFromUserSessionID(SessionID SessID) (*Lobby, error) {
  query := bson.M{
    "users": bson.M{
      "$elemMatch": bson.M{ "sessionID": SessionID, },
    },
  }

  var lobby Lobby
  result := LobbyDB.Coll.FindOne(LobbyDB.Context, query)
  if err := result.Err(); err != nil {
    return nil, err
  }
  if err := result.Decode(&lobby); err != nil {
    return nil, err
  }
  return &lobby, nil
}

func SaveLobby(lobby *Lobby) error {
  _, err := LobbyDB.Coll.InsertOne(LobbyDB.Context, lobby)
  return err
}

