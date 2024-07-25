package GdHandler

import (
  "errors"
  "reflect"
  "strconv"
  "encoding/json"
  
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)


/// The lobby handling function responsible for connecting players to their respective lobby
func (lobby *Lobby) wsHandler(c *gin.Context) error {
  var myIndex int = -1
  conn, err := lobby.Upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    return errors.New("wsHandler: Upgrade error\n" + err.Error())
  }
  defer func () {
    lobby.PlayerMutex.Lock()
    lobby.disconnectPlayer(myIndex)
    lobby.PlayerMutex.Unlock()
  }()

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
      if conn.WriteMessage(websocket.BinaryMessage, []byte{0x00, byte(myIndex)}) == nil {
        break
      }
    } else {
      _ = conn.WriteMessage(websocket.BinaryMessage, append([]byte{0xff}, []byte(err.Error())...))
    }
  }

  lobby.PlayerMutex.Lock()
  err = lobby.connectPlayer(conn, myIndex)
  lobby.PlayerMutex.Unlock()
  // Handle outbound data
  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      return errors.New("wsHandler: Read error:" + err.Error())
    } else if messageType == websocket.CloseMessage { break }

    // Async can cause UB as values can be modified while another goroutine is in flight
    if messageType == websocket.TextMessage {
      err = lobby.handleTextMessage(message, conn)
    } else {
      err = errors.New("wsHandler: Invalid messageType: " + strconv.Itoa(messageType))
    }

    if err != nil {
      val, jsonerr := json.Marshal(struct {Error string}{err.Error()})
      if jsonerr == nil {
        // log.Println("Text message marshal error: ", err)
        continue
      }
      _ = conn.WriteMessage(websocket.BinaryMessage, val)
      // log.Println("wsHandler: error:", err)
    }
  }
  return nil
}

/// Checks user auth
func (lobby *Lobby) getUserAuth(messageType int, message []byte) (int, error) {
  if (messageType != websocket.BinaryMessage) {
    return -1, errors.New("getUserAuth: Invalid messageType")
  }

  if len(message) != 64 {
    return -1, errors.New("getUserAuth: Invalid token length!")
  }

  // Validate token
  for i, player := range lobby.Players {
    if reflect.DeepEqual(*player.SessionID, [64]byte(message)) {
      return i, nil
    }
  }

  return -1, errors.New("getUserAuth: Denied")
}

func (lobby *Lobby) connectPlayer(conn *websocket.Conn, myIndex int) error {
  if myIndex < 0 { 
    return errors.New("lobby.connectPlayer: myIndex underflow")
  }
  if len(lobby.Players) >= myIndex {
    return errors.New("lobby.connectPlayer: myIndex overflow")
  }

  lobby.Players[myIndex].Conn = conn
  return nil
}

func (lobby *Lobby) disconnectPlayer(myIndex int) {
  if myIndex < 0 { return }
  if len(lobby.Players) >= myIndex { return }

  if conn := lobby.Players[myIndex].Conn; conn != nil {
    conn.Close()
  }

  lobby.Players[myIndex].Conn = nil
}

