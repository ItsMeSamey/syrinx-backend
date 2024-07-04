package GdHandler

import (
  "errors"
  "log"
  "net/http"
  "strconv"
  "sync"

  "ccs.ctf/DB"

  "github.com/gorilla/websocket"
)

type Players struct {
  id int
  conn *websocket.Conn
}

type Lobby struct {
  players []Players
  dataPool []byte
  dataMutex sync.Mutex
  playerMutex sync.Mutex
  upgrader websocket.Upgrader
}

func MakeLobby(lobbyID int) *Lobby {
  return &Lobby {
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

func (lobby *Lobby) wsHandler(w http.ResponseWriter, r *http.Request) error {
  conn, err := lobby.upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Print("wsHandler: Upgrade error:", err)
    return err
  }

  var user *DB.User = nil

  defer conn.Close()

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
    if user == nil {
      user, err = getUserAuth(messageType, message)
      if err == nil {
        _ = conn.WriteMessage(messageType, []byte("0Success"))
      } else {
        _ = conn.WriteMessage(messageType, []byte("1Authentication Error"))
      }
      continue
    }

    if messageType == websocket.TextMessage {
      err = lobby.handleTextMessage(message)
    } else if messageType == websocket.BinaryMessage{
      // Consider making this async
      // This may lead to partially unmerged packets as batching may happen before
      // merging is completed by all of the players, but is that even a problem?
      lobby.handleBinaryMessage(message)
    } else {
      err = errors.New("wsHandler: Invalid messageType: " + strconv.Itoa(messageType))
    }
    if err != nil {
      log.Println("wsHandler: error:", err)
      continue
    }

    log.Printf("wsHandler: Received message: %s\n", message)

    // Optionally send a response message
    err = conn.WriteMessage(messageType, []byte("Hello from server!"))
    if err != nil {
      log.Println("wsHandler: Write error:", err)
      continue
    }
  }
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
      // Also set a context timeout to prevent a deadlock
      _ = i.conn.WriteMessage(websocket.BinaryMessage, lobby.dataPool)
      wg.Done()
    })()
  }
  wg.Wait()
}

