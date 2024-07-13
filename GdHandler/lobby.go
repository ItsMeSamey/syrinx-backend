package GdHandler

import (
  "errors"
  "log"
  "strconv"
  "sync"

  "ccs.ctf/DB"
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)

type Lobby struct {
  ID          DB.ObjID
  players     []DB.Player
  playercount byte
  playerMutex sync.RWMutex
  upgrader    websocket.Upgrader
  deadtime    byte
}

/// The lobby handling function responsible for connecting players to their respective lobby
func (lobby *Lobby) wsHandler(c *gin.Context) {
  conn, err := lobby.upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    log.Print("wsHandler: Upgrade error:", err)
  }
  defer func () {
    lobby.playerMutex.Lock()
    lobby.playercount -= 1
    lobby.playerMutex.Unlock()
    conn.Close()
  }()

  var myIndex byte
  // Authanticate the user
  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("wsHandler: Read error:", err)
      continue
    } else if messageType == websocket.CloseMessage {
      _ = conn.Close()
      return
    }
    myIndex, err = lobby.getUserAuth(messageType, message)
    if err == nil {
      if conn.WriteMessage(websocket.BinaryMessage, []byte{0, myIndex}) == nil {
        break
      }
    } else {
      _ = conn.WriteMessage(websocket.BinaryMessage, []byte{0})
    }
  }

  // Create a player receiving channel
  channel := make(chan []byte, 128)
  lobby.playerMutex.Lock()
  lobby.players[myIndex].IN = channel
  lobby.playerMutex.Unlock()

  // Delete the player receiving channel in the end
  defer func () {
    lobby.playerMutex.Lock()
    lobby.players[myIndex].IN = nil
    lobby.playerMutex.Unlock()
    close(channel)
  }()

  // Handle incoming data
  go func () {
    for packet := range channel {
      _ = conn.WriteMessage(websocket.BinaryMessage, packet)
    }
  }()

  // Handle outbound data
  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("wsHandler: Read error:", err)
      continue
    }
    if messageType == websocket.CloseMessage {
      break
    }

    // Async can cause UB as values can be modified while another goroutine is in fligt
    if messageType == websocket.TextMessage {
      err = lobby.handleTextMessage(myIndex, message)
    } else if messageType == websocket.BinaryMessage{
      err = lobby.handleBinaryMessage(myIndex, message)
    } else {
      err = errors.New("wsHandler: Invalid messageType: " + strconv.Itoa(messageType))
    }
    if err != nil {
      log.Println("wsHandler: error:", err)
      continue
    }
  }
}

