package DB

import "errors"

type Lobby struct {
  ID   ObjID    `bson:"_id,omitempty"`
  UserIDs []ObjID `bson:"users"`
}

func GetLobby(lobbyID ObjID) ([]ObjID, error) {
  var lobby Lobby
  err := LobbyDB.get("_id", lobbyID, &lobby)
  if err != nil {
    return nil, err
  }
  if len(lobby.UserIDs) == 0 {
    return nil, errors.New("GetLobby: Lobby is empty")
  }
  return lobby.UserIDs, nil
}

