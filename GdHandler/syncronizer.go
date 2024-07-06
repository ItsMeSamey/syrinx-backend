package GdHandler

import (
  "errors"
  "net/http"
  "strconv"
  "sync"
  "time"
  
  "ccs.ctf/DB"
  
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)

var Lobbies map[int]*Lobby = make(map[int]*Lobby)

func LobbyHandler(gc *gin.Context) error {
  idStr := gc.Param("lobbyID")
  lobbyID, err := strconv.Atoi(idStr)
  if err != nil {
    gc.AbortWithStatus(http.StatusNotFound)
    return errors.New("LobbyHandler: lobbyID is not an int")
  }

  lobby, ok := Lobbies[lobbyID]
  if ok {
    lobby.wsHandler(gc)
  } else {
    exist, err := DB.DoesExistLobby([]byte(idStr))
    if err != nil {
      gc.AbortWithStatus(http.StatusInternalServerError)
      return err
    }
    if exist {
      lobby = makeLobby(lobbyID)
      Lobbies[lobbyID] = lobby
      lobby.wsHandler(gc)
      Syncronizer(lobby)
    } else {
      gc.AbortWithStatus(http.StatusNotFound)
      return errors.New("LobbyHandler: lobby does not exist")
    }
  }
  return nil
}

func Syncronizer(lobby *Lobby) {
start:
  lobby.dataMutex.Lock()
  defer lobby.dataMutex.Unlock()
  lobby.playerMutex.Lock()
  defer lobby.playerMutex.Unlock()
  playerCount := len(lobby.players)
  if playerCount == 0 {
    delete(Lobbies, lobby.lobbyID)
    return
  }

  var wg sync.WaitGroup
  wg.Add(playerCount)
  for _, i := range lobby.players {
    go (func ()  {
      // HACK: Donot discard the error try again
      // FIXME: Set a timeout to prevent a deadlock (maybe using context)
      _ = i.conn.WriteMessage(websocket.BinaryMessage, lobby.dataPool)
      wg.Done()
    })()
  }
  wg.Wait()
  time.Sleep(50*time.Millisecond)
  goto start
}

