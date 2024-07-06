package DB

import (

)

type Lobby struct {
  ID   string    `bson:"_id,omitempty"`
  Teams []string `bson:"teams"`
}

func UserInLobby(_id string) (bool, error) {
  return LobbyDB.exists("_id", _id)
}

func HasUser(lobbyID, userID string) (bool, error) {
  var lobby Lobby
  err := LobbyDB.get("_id", lobbyID, &lobby)
  if err != nil {
    return false, err
  }
  for _, id := range lobby.Teams {
    if userID == id {
      return true, nil
    }
  }
  return false, nil
}

