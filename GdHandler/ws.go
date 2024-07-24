package GdHandler

import (
  "log"
  "errors"
  "reflect"
  "strconv"
  "net/http"
  "encoding/json"
  
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)


/// origin checker for websocket connections
func originChecker(r *http.Request) bool {
  return true
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
      if jsonerr == nil { log.Println("Text message marshal error: ", err) }
      _ = conn.WriteMessage(websocket.BinaryMessage, val)
      log.Println("wsHandler: error:", err)
    }
  }
  return nil
}

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

/// Forcefully close the lobby
func (lobby *Lobby) delete() {
  lobby.deleteAllPlayer()

  lobbiesMutex.Lock()
  defer lobbiesMutex.Unlock()

  // needs to be called again after mutex locking
  lobby.deleteAllPlayer()

  if _, ok := lobbies[*(lobby.Team.TeamID)]; ok {
    delete(lobbies, *(lobby.Team.TeamID))
  }
}

func (lobby *Lobby) deleteAllPlayer() {
  lobby.PlayerMutex.Lock()
  defer lobby.PlayerMutex.Unlock()

  for i := range lobby.Players {
    if lobby.Players[i].Conn != nil {
      lobby.Players[i].Conn.Close()
      lobby.Players[i].Conn = nil
    }
  }
}

