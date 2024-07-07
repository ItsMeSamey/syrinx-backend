package GdHandler

import (
  "errors"
  "log"
  "strconv"
  "sync"

  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
)

type Lobby struct {
  ID      string
  players []*Player
  dataPool []byte
  dataMutex sync.Mutex
  playerMutex sync.Mutex
  upgrader websocket.Upgrader
}

func makeLobby(ID string) *Lobby {
  return &Lobby {
    ID: ID,
    players: make([]*Player, 0, 32),
    dataPool: make([]byte, 0, 1024),
    dataMutex: sync.Mutex{},
    playerMutex: sync.Mutex{},
    upgrader: websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024 },
  }
}

func (lobby *Lobby) handleTextMessage(message []byte) error {
  // TODO: implement
  _ = message
  return nil
}

func (lobby *Lobby) handleBinaryMessage(message []byte) {
  lobby.dataMutex.Lock()
  defer lobby.dataMutex.Unlock()
  lobby.dataPool = append(lobby.dataPool, message...)
}

func (lobby *Lobby) wsHandler(gc *gin.Context, index int) {
  conn, err := lobby.upgrader.Upgrade(gc.Writer, gc.Request, nil)
  if err != nil {
    log.Print("wsHandler: Upgrade error:", err)
  }
  defer conn.Close()

  myself := lobby.players[index]
  myself.WS = conn
  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("wsHandler: Read error:", err)
      continue
    }
    if messageType == websocket.CloseMessage {
      _ = conn.Close()
      return
    }
    myself.ID, err = lobby.getUserAuth(messageType, message)
    if err == nil && myself.ID != "" {
      if conn.WriteMessage(messageType, []byte("0Success")) == nil {
        break
      }
    } else {
      _ = conn.WriteMessage(messageType, []byte("1Authentication Error"))
    }
  }

  for {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("wsHandler: Read error:", err)
      continue
    }
    if messageType == websocket.CloseMessage {
      conn.Close()
      break
    }

    // If handlers are made async, programme will eventually lock up on slow pc's
    if messageType == websocket.TextMessage {
      err = lobby.handleTextMessage(message)
    } else if messageType == websocket.BinaryMessage{
      lobby.handleBinaryMessage(message)
    } else {
      err = errors.New("wsHandler: Invalid messageType: " + strconv.Itoa(messageType))
    }
    if err != nil {
      log.Println("wsHandler: error:", err)
      continue
    }
  }

  lobby.playerMutex.Lock()
  defer lobby.playerMutex.Unlock()
  if len(lobby.players) == index+1 {
    lobby.players = lobby.players[:index-1]
    return
  } 
  lobby.players = append(lobby.players[:index], lobby.players[index+1:]...)
  return
}

