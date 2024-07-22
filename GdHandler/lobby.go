package GdHandler

import (
  "sync"
  "errors"

  "ccs.ctf/DB"

  "github.com/gorilla/websocket"
)

type Lobby struct {
  Lobby       *DB.Lobby
  Playercount byte
  PlayerMutex sync.RWMutex
  Upgrader    websocket.Upgrader
  Deadtime    byte
}

/// Function to make an empty lobby struct
func makeLobby(ID DB.ObjID) (*Lobby, error) {
  lobby, err := DB.LobbyFromID(ID)
  if err != nil {
    return nil, err
  }

  if err := lobby.PopulateTeams(); err != nil {
    return nil, errors.New("makeLobby: Lobby.populateTeams error\n" + err.Error())
  }

  return &Lobby {
    Lobby: lobby,
    Playercount: 0,
    PlayerMutex: sync.RWMutex{},
    Upgrader: websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024, CheckOrigin: originChecker},
    Deadtime: 0,
  }, nil
}

func LobbyIDFromUserSessionID(SessionID DB.SessID) (DB.ObjID, error) {
  lobby, exists, err := DB.LobbyFromUserSessionID(SessionID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.LobbyFromUserSessionID error\n" + err.Error())
  }
  
  if exists {
    return lobby.ID, nil
  }

  template, err := DB.NewLobbyTemplate(SessionID)
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.NewLobbyTemplate error\n" + err.Error())
  }
  
  lobby, exists, err = DB.GetIncompleteLobby()
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: DB.GetIncompleteLobby error\n" + err.Error())
  }

  insertable := &Lobby {
    Lobby: nil,
    Playercount: 0,
    PlayerMutex: sync.RWMutex{},
    Upgrader: websocket.Upgrader{ ReadBufferSize:  1024, WriteBufferSize: 1024, CheckOrigin: originChecker},
    Deadtime: 0,
  }

  if exists {
    // A partially filled lobby found
    err = lobby.PopulateTeams()
    if err != nil {
      return nil, errors.New("LobbyIDFromUserSessionID: Lobby.PopulateTeams error\n" + err.Error())
    }

    lobby.Merge(template)
    err = lobby.Sync(5)
    if err != nil {
      return nil, errors.New("LobbyIDFromUserSessionID: Lobby.Sync error\n" + err.Error())
    }

    insertable.Lobby = lobby
    lobbiesMutex.Lock()
    val, ok := lobbies[*(lobby.ID)]
    if ok {
      // Lobby already active, updating
      val.Lobby = lobby
    } else {
      // Make lobby active
      lobbies[*(lobby.ID)] = insertable
    }
    lobbiesMutex.Unlock()

    if !ok {
      go watchdog(insertable)
    }
    return lobby.ID, nil
  }

  // Creating a new lobby
  err = template.CreateInsert()
  if err != nil {
    return nil, errors.New("LobbyIDFromUserSessionID: Lobby.CreateInsert error\n" + err.Error())
  }

  insertable.Lobby = template
  // Activating the new lobby
  lobbiesMutex.Lock()
  lobbies[*(template.ID)] = insertable
  lobbiesMutex.Unlock()

  go watchdog(insertable)
  return template.ID, nil
}

