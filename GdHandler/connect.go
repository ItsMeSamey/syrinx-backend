package GdHandler

import (
  "errors"
  
  "ccs.ctf/DB"

  "github.com/gin-gonic/gin"
)

//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Connect to a lobby if one exists or add it to the lobby map
func getAddedLobby(ID DB.TID, execFunc func(*Lobby)error) (*Lobby, error) {
  start:
  lobbiesMutex.RLock()
  val, ok := lobbies[*ID]
  lobbiesMutex.RUnlock()
  if ok {
    val.PlayerMutex.Lock()
    if val.Deadtime >= 10 { goto start }
    val.Playercount += 1
    val.PlayerMutex.Unlock()
    return val, execFunc(val)
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
        return nil, errors.New("ConnectToLobby: error while lobby creation\n" + err.Error())
      }
      lobbies[*ID] = val
      val.Playercount += 01
      lobbiesMutex.Unlock()
      
      go func(){
        TEAMSMutex.Lock()
        TEAMS = append(TEAMS, val.Team)
        TEAMSMutex.Unlock()
      }()

      return val, execFunc(val)
    }
  }
}

func ConnectToLobby(ID DB.TID, c *gin.Context) error {
  _, err := getAddedLobby(ID, func(lobby *Lobby) error { return lobby.wsHandler(c)})
  return err
}

