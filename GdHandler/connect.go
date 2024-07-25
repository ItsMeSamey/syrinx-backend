package GdHandler

import (
  "errors"
  
  "ccs.ctf/DB"

  "github.com/gin-gonic/gin"
)

//! WARNING: DO NOT MESS WITH THIS UNLESS YOU KNOW WHAT YOU ARE DOING
/// Connect to a lobby if one exists or add it to the lobby map
func getAddedLobby(ID DB.TID) (*Lobby, error) {
  start:
  lobbiesMutex.RLock()
  val, ok := lobbies[*ID]
  lobbiesMutex.RUnlock()
  if ok {
    return val, nil
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
      lobbiesMutex.Unlock()
      
      go func(){
        TEAMSMutex.Lock()
        TEAMS = append(TEAMS, val.Team)
        TEAMSMutex.Unlock()
      }()

      return val, nil
    }
  }
}

func ConnectToLobby(ID DB.TID, c *gin.Context) error {
  val, err := getAddedLobby(ID)
  if err != nil {
    return err
  }
  return val.wsHandler(c)
}

