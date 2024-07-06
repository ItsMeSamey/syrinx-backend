package DB

import (

)

type Lobby struct {
  ID   string    `bson:"_id,omitempty"`
  Teams []string `bson:"teams"`
}

func UserInLobby(id []byte) (bool, error) {
  return LobbyDB.DoesExist(lobbyBucket, id)
}

func (lobby *Lobby) HasUser(username string) (bool, error) {
  // Implement
  return false, nil
}



