package GdHandler

import (
	"net/http"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MAX_TRIES = 5

var (
  lobbies map[primitive.ObjectID]*Lobby = make(map[primitive.ObjectID]*Lobby)
  lobbiesMutex sync.RWMutex = sync.RWMutex{}
)

func originChecker(r *http.Request) bool {
  return true
}
