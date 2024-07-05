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

type Players struct {
  DB.User
  conn *websocket.Conn
}

type Lobby struct {
  ID int
  players []Players
  dataPool []byte
  dataMutex sync.Mutex
  playerMutex sync.Mutex
  upgrader websocket.Upgrader
}

func MakeLobby(lobbyID int) *Lobby {
  return &Lobby {
    ID: lobbyID,
    players: make([]Players, 0, 32),
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

func (lobby *Lobby) wsHandler(gc *gin.Context) error {
  conn, err := lobby.upgrader.Upgrade(gc.Writer, gc.Request, nil)
  if err != nil {
    log.Print("wsHandler: Upgrade error:", err)
    return err
  }
  defer conn.Close()

  var user *DB.User = nil
  for user == nil {
    messageType, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("wsHandler: Read error:", err)
      continue
    }
    if messageType == websocket.CloseMessage {
      conn.Close()
      break
    }
    user, err = lobby.getUserAuth(messageType, message)
    if err == nil && user != nil {
      _ = conn.WriteMessage(messageType, []byte("0Success"))
    } else {
      _ = conn.WriteMessage(messageType, []byte("1Authentication Error"))
    }
    continue
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

  // Remove the user from the lobby and `players` before exit
  return nil
}

func (lobby *Lobby) Syncronizer() {
  lobby.dataMutex.Lock()
  defer lobby.dataMutex.Unlock()
  lobby.playerMutex.Lock()
  defer lobby.playerMutex.Unlock()

  var wg sync.WaitGroup
  wg.Add(len(lobby.players))
  for _, i := range lobby.players {
    go (func ()  {
      // HACK: Donot discard the error try again
      // FIXME: Set a timeout to prevent a deadlock (maybe using context)
      _ = i.conn.WriteMessage(websocket.BinaryMessage, lobby.dataPool)
      wg.Done()
    })()
  }
  wg.Wait()
}

