package GdHandler

import (
  "sync"
  
  "ccs.ctf/DB"
)

const MAX_TRIES = 5

var (
  /// The active lobbies are stored here
  lobbies map[[3]byte]*Lobby = make(map[[3]byte]*Lobby)
  lobbiesMutex sync.RWMutex = sync.RWMutex{}

  /// Set Level of the lobbies
  LEVEL = DB.State.Level
)

func Init() error {
  DB.Callbacks["level updater"] = func (prev, cur *DB.STATE) {
    if prev.Level != cur.Level {
      LEVEL = cur.Level
    }
  }

  return nil
}

