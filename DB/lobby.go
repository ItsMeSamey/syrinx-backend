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

/// Synchronously return the functions
func ForEachLobby(fn func (k, v []byte) error) error {
  return LobbyDB.forEachInBucket(lobbyBucket, fn)
}

func (lobby *Lobby) HasUser(username string) (bool, error) {
  // Implement
  return false, nil
}


