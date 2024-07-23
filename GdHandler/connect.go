package GdHandler

import (
  "time"
  "errors"
  
  "ccs.ctf/DB"

  "github.com/gin-gonic/gin"
)

//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Connect to a lobby if one exists or add it to the lobby map
func ConnectToLobby(ID DB.TID, c *gin.Context) error {
  start:
  lobbiesMutex.RLock()
  val, ok := lobbies[*ID]
  lobbiesMutex.RUnlock()
  if ok {
    val.PlayerMutex.Lock()
    if val.Deadtime >= 10 { goto start }
    val.Playercount += 1
    val.PlayerMutex.Unlock()
    return val.wsHandler(c)
  } else {
    // A user will be stranded in a isolated lobby if thisis ignored
    lobbiesMutex.Lock()
    val, ok := lobbies[*ID]
    
    if ok {// If a lobby was created when we switched locks !!
      lobbiesMutex.Unlock()
      goto start
    } else {
      var err error
      val, err = makeLobby(ID)
      if err != nil {
        lobbiesMutex.Unlock()
        return errors.New("ConnectToLobby: error while lobby creation\n" + err.Error())
      }
      lobbies[*ID] = val
      lobbiesMutex.Unlock()

      val.Playercount += 1
      go watchdog(val)
      return val.wsHandler(c)
    }
  }
}

//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Automatically close the lobby when there is no one in it
/// This function bocks (a long as lobby exists),
/// and probably should be async
func watchdog(lobby *Lobby) {
  lobby.Deadtime = 0
  sleepTime := (30 + (lobby.Team.TeamID[0]&31))
begin:
  time.Sleep(time.Duration(sleepTime) * time.Second)
  lobby.PlayerMutex.RLock()
  if lobby.Playercount == 0 {
    lobby.Deadtime += 1
  }
  lobby.PlayerMutex.RUnlock()

  if lobby.Deadtime >= 10 { // Lobby timeout 5~10 minutes
    lobbiesMutex.Lock()
    lobby.PlayerMutex.Lock()
    if lobby.Playercount == 0 {
      // Delete the lobby
      delete(lobbies, *lobby.Team.TeamID);
      lobby.PlayerMutex.Unlock()
      lobbiesMutex.Unlock()
      return
    }
    lobby.Deadtime = 0
    lobby.PlayerMutex.Unlock()
    lobbiesMutex.Unlock()
  }
  goto begin
}

