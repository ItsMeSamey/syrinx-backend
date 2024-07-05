package DB

import (

)

const (
  lobbyBucket = "lobby"
)

type Lobby struct {
  lobbyID int
	Teams []string
}

func DoesExistLobby(id []byte) (bool, error) {
  return LobbyDB.DoesExist(lobbyBucket, id)
}

func (lobby *Lobby) HasUser(username string) (bool, error) {
  // Implement
  return false, nil
}



