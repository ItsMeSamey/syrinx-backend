package GdHandler

import (
  "sort"
  "sync"
  "time"
  "errors"
  
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
  DB.Callbacks["level updater"] = func (prev, cur *DB.STATE) {
    if prev.Level != cur.Level || prev.Keep != cur.Keep {
      LEVEL = cur.Level
      go changeLevelTo(cur.Level, cur.Keep)
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

  for _, team := range TEAMS {
    lobbies[*(team.TeamID)] = makeLobbyFromTeam(team)
  }

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

