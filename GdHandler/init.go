package GdHandler

import (
	"sync"

	"ccs.ctf/DB"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var lobbies map[primitive.ObjectID]*Lobby
var lobbiesMutex sync.RWMutex = sync.RWMutex{}

/// Connect to a lobby if one exists or add it to the lobby map
func ConnectToLobby(lobby *DB.Lobby, c *gin.Context) {
	/// WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
	lobbiesMutex.RLock()
	val, ok := lobbies[*lobby.ID]
	if ok {
		lobbiesMutex.RUnlock()
		val.wsHandler(c)
	} else {
		lobbiesMutex.RUnlock()

		/// Check if a lobby was created when we switched locks !!
		/// A user will be stranded in a isolated lobby if this happens
		lobbiesMutex.Lock()
		val, ok := lobbies[*lobby.ID]
		if ok {
			lobbiesMutex.Unlock()
			val.wsHandler(c)
		} else {
			val = makeLobby(lobby)
			lobbies[*lobby.ID] = val
			lobbiesMutex.Unlock()
			// create a timeout function
			val.wsHandler(c)
		}
	}
}

