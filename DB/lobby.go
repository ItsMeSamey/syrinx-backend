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

func ForEachLobby(fn func (k, v []byte) error) error {
  return LobbyDB.forEachInBucket(lobbyBucket, fn)
}

func (lobby *Lobby) HasUser(username string) (bool, error) {
	// Implement
	return false, nil
}


