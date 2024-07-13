package GdHandler

import (
  "sync"
  "time"
  
  "ccs.ctf/DB"
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
  lobbies map[primitive.ObjectID]*Lobby
  lobbiesMutex sync.RWMutex = sync.RWMutex{}
)

/// Function to make an empty lobby struct
func makeLobby(lobby *DB.Lobby) *Lobby {
  return &Lobby {
    ID: lobby.ID,
    players: lobby.Players,
    playercount: 0,
    playerMutex: sync.RWMutex{},
    upgrader: websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024 },
    deadtime: 0,
  }
}


//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Automatically close the lobby when there is no one in it
/// This function bocks (a long as lobby exists),
/// and probably should be async
func watchdog(lobby *Lobby) {
  lobby.deadtime = 0
  sleepTime := (30 + (lobby.ID[0]&31))
begin:
  time.Sleep(time.Duration(sleepTime) * time.Second)
  lobby.playerMutex.RLock()
  if lobby.playercount == 0 {
    lobby.deadtime += 1
  }
  lobby.playerMutex.RUnlock()

  if lobby.deadtime >= 10 { // Lobby timeout 5~10 minutes
    lobbiesMutex.Lock()
    lobby.playerMutex.Lock()
    if lobby.playercount == 0 {
      // Delete the lobby
      delete(lobbies, *lobby.ID);
      lobby.playerMutex.Unlock()
      lobbiesMutex.Unlock()
      return
    }
    lobby.deadtime = 0
    lobby.playerMutex.Unlock()
    lobbiesMutex.Unlock()
  }
  goto begin
}

//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Connect to a lobby if one exists or add it to the lobby map
func ConnectToLobby(lobby *DB.Lobby, c *gin.Context) {
  start:
  lobbiesMutex.RLock()
  val, ok := lobbies[*lobby.ID]
  lobbiesMutex.RUnlock()
  if ok {
    val.playerMutex.Lock()
    if val.deadtime >= 10 { goto start }
    val.playercount += 1
    go val.wsHandler(c)
  } else {
    // A user will be stranded in a isolated lobby if thisis ignored
    lobbiesMutex.Lock()
    val, ok := lobbies[*lobby.ID]
    
    if ok {// If a lobby was created when we switched locks !!
      lobbiesMutex.Unlock()
      goto start
    } else {
      val = makeLobby(lobby)
      lobbies[*lobby.ID] = val
      lobbiesMutex.Unlock()
      val.playercount += 1
      go val.wsHandler(c)
      go watchdog(val)
    }
  }
}

