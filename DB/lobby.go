package DB

import "errors"

type Lobby struct {
  ID   string    `bson:"_id,omitempty"`
  UserIDs []string `bson:"users"`
}

func GetLobby(lobbyID string) ([]string, error) {
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

