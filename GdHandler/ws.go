package GdHandler

import (
  "log"
  "errors"
  "strconv"
  "reflect"

  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)

/// Checks user auth
func (lobby *Lobby) getUserAuth(messageType int, message []byte) (byte, error) {
  if (messageType != websocket.BinaryMessage) {
    return 0, errors.New("getUserAuth: Invalid messageType")
  }

  if len(message) != 64 {
    return 0, errors.New("getUserAuth: Invalid token length!")
  }

  // Validate token
  for i, player := range lobby.Players {
    if reflect.DeepEqual(*player.SessionID, [64]byte(message)) {
      return byte(i), nil
    }
  }

  return 0, errors.New("getUserAuth: Denied")
}

/// The lobby handling function responsible for connecting players to their respective lobby
func (lobby *Lobby) wsHandler(c *gin.Context) error {
  conn, err := lobby.Upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    return errors.New("wsHandler: Upgrade error\n" + err.Error())
  }
  defer func () {
    lobby.PlayerMutex.Lock()
    lobby.Playercount -= 1
    lobby.PlayerMutex.Unlock()
    conn.Close()
  }()

  var myIndex byte
  // Authanticate the user
  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      return errors.New("wsHandler: Read error:" + err.Error())
    } else if messageType == websocket.CloseMessage {
      _ = conn.Close()
      return errors.New("wsHandler: Connection closed without auth")
    }
    myIndex, err = lobby.getUserAuth(messageType, message)
    if err == nil {
      if conn.WriteMessage(websocket.BinaryMessage, []byte{0x00, myIndex}) == nil {
        break
      }
    } else {
      _ = conn.WriteMessage(websocket.BinaryMessage, append([]byte{0xff}, []byte(err.Error())...))
    }
  }

  // Create a player receiving channel
  channel := make(chan []byte, 128)
  lobby.PlayerMutex.Lock()
  if lobby.Players[myIndex].IN != nil {
    close(lobby.Players[myIndex].IN)
  }
  lobby.Players[myIndex].IN = channel
  lobby.PlayerMutex.Unlock()

  // Delete the player receiving channel in the end
  defer func () {
    lobby.PlayerMutex.Lock()
    if lobby.Players[myIndex].IN != nil {
      close(channel)
      lobby.Players[myIndex].IN = nil
    }
    lobby.PlayerMutex.Unlock()
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
      return errors.New("wsHandler: Read error:" + err.Error())
    }
    if messageType == websocket.CloseMessage {
      break
    }

    // Async can cause UB as values can be modified while another goroutine is in flight
    if messageType == websocket.TextMessage {
      err = lobby.handleTextMessage(myIndex, message, conn)
    } else if messageType == websocket.BinaryMessage{
      err = lobby.handleBinaryMessage(myIndex, message)
    } else {
      err = errors.New("wsHandler: Invalid messageType: " + strconv.Itoa(messageType))
    }
    if err != nil {
      log.Println("wsHandler: error:", err)
      _ = conn.WriteMessage(websocket.BinaryMessage, append([]byte{0xff}, []byte(err.Error())...))
      continue
    }
  }
  return nil
}

