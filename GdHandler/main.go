package GdHandler

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sort"
	"sync"
	"time"

	"ccs.ctf/DB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MAX_TRIES = 5

var (
  /// The active lobbies are stored here
  lobbies map[[3]byte]*Lobby = make(map[[3]byte]*Lobby)
  lobbiesMutex sync.RWMutex = sync.RWMutex{}

  /// Set Level of the lobbies
  LEVEL = DB.State.Level
  /// All teams array
  TEAMS []*DB.Team
  TEAMSMutex sync.RWMutex = sync.RWMutex{}
)

func Init() error {
  DB.Callback = func (prev, cur *DB.STATE) {
    if prev.Level != cur.Level || prev.Keep != cur.Keep {
      LEVEL = cur.Level
      keep := cur.Keep
      go changeLevelTo(LEVEL, keep)
    }

    if cur.TeamExceptions != nil {
      exp := cur.TeamExceptions
      go syncTeams(exp)
    }
  }

  batchSize := int32(1024)
  cursor, err :=  DB.TeamDB.Coll.Find(DB.TeamDB.Context, bson.M{}, &options.FindOptions{
    BatchSize: &batchSize,
    Sort: bson.M{"points": -1},
  })
  if err != nil {
    return errors.New("GdHandler.Init: error in DB.Find\n" + err.Error())
  }

  if err = cursor.All(DB.TeamDB.Context, &TEAMS); err != nil {
    return errors.New("GdHandler.Init: error in cursor.All\n" + err.Error())
  }

  var wg sync.WaitGroup

  for _, team := range TEAMS {
    wg.Add(1)
    go func () {
      defer wg.Done()
      lobby := makeLobbyFromTeam(team)
      start:
      err := lobby.populatePlayers()
      if err != nil {
        goto start
      }
      lobbies[*(team.TeamID)] = lobby
    }()
  }

  wg.Wait()

  go func ()  {
    for {
      TEAMSMutex.Lock()
      sortTeams()
      TEAMSMutex.Unlock()
      time.Sleep(2500 * time.Millisecond)
    }
  }()

  return nil
}

func changeLevelTo(level, keep int) {
  TEAMSMutex.Lock()
  lobbiesMutex.RLock()
  for _, lobby := range lobbies {
    lobby.delete()
  }
  lobbiesMutex.RUnlock()

  sortTeams()

  for i, team := range TEAMS {
    var final int = 0
    if i <= keep {
      final = level
    } else {
      final = level-1
    }
    
    if team.Level != final {
      team.Level = final
      go team.Sync(5)
    }
  }

  TEAMSMutex.Unlock()
}

func sortTeams() {
  sort.Slice(TEAMS, func (i, j int) bool {
    ti := TEAMS[i]
    tj := TEAMS[j]
    if ti.Points != tj.Points {
      return ti.Points > tj.Points
    }

    timei := teamTimeSum(ti)
    timej := teamTimeSum(tj)
    if timei != timej {
      return timei < timej
    }

    return ti.TeamName > tj.TeamName
  })
}
func teamTimeSum(team *DB.Team) int64 {
  sum := int64(0)
  for _, val := range team.Solved {
    sum += val
  }
  return sum
}

func syncTeams(exceptions []any) {
  if len(exceptions) == 0 { return }
  if len(exceptions) == 1 {
    val, ok := exceptions[0].(string)
    if ok {
      if val == "" {
        lobbiesMutex.Lock()
        for _, val := range lobbies {
          resyncLobby(val)
        }
      }
    }
  }
  for _, idVar := range exceptions {
    var id [3]byte
    switch idTyped := idVar.(type) {
    case []byte:
      if len(idTyped) != 3 { continue }
      id = [3]byte(idTyped)

    case string:
      if len(idTyped) == 6 {
        data, err := hex.DecodeString(idTyped)
        if err != nil { continue }
        if len(data) != 3 { continue }
        id = [3]byte(data)
      } else if len(idTyped) == 4 {
        data, err := base64.StdEncoding.DecodeString(idTyped)
        if err != nil { continue }
        if len(data) != 3 { continue }
        id = [3]byte(data)
      } else { continue }

    default:
      continue
    }

    val, ok := lobbies[id]
    if ok { go resyncLobby(val) }
  }
}
func resyncLobby(val *Lobby) {
  team, err := DB.TeamByTeamID(val.Team.TeamID)
  lobby := makeLobbyFromTeam(team)

  if err != nil { return }
  err = lobby.populatePlayers()

  if err != nil { return }
  val.Team = team
  val.PlayerMutex.Lock()
  val.delete()
  val.Players = lobby.Players
  val.PlayerMutex.Unlock()
}

