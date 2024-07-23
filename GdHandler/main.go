package GdHandler

import (
  "sync"
  "net/http"
)

const MAX_TRIES = 5

var (
  /// The active lobbies are stored here
  lobbies map[[3]byte]*Lobby = make(map[[3]byte]*Lobby)
  lobbiesMutex sync.RWMutex = sync.RWMutex{}
)

func originChecker(r *http.Request) bool {
  return true
}
