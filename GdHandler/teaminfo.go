package GdHandler

import (
	"errors"

	"ccs.ctf/DB"

	"go.mongodb.org/mongo-driver/bson"
)

type TeamMate struct {
  Username  string `bson:"user"`
  Email     string `bson:"mail"`
  DiscordID string `bson:"discordID"`
}

func GetTeamAndPlayers(ID DB.TID) ([]TeamMate, *DB.Team, error) {
  var team *DB.Team
  var all []TeamMate
  var err error

  lobbiesMutex.RLock()
  defer lobbiesMutex.RUnlock()
  val, ok := lobbies[*ID]
  if ok {
    team = val.Team
    for _, player := range val.Players {
      all = append(all, TeamMate{
        Username:  player.Username,
        Email:     player.Email,
        DiscordID: player.DiscordID,
      })
    }
  } else {
    team, err = DB.TeamByTeamID(ID)
    if err != nil {
      return nil, nil, errors.New("GetTeamAndPlayers: Error getting team\n" + err.Error())
    }

    cursor, err := DB.UserDB.Coll.Find(DB.UserDB.Context, bson.M{"teamID": ID})

    err = cursor.All(DB.UserDB.Context, &all)
    if err != nil {
      return nil, nil, errors.New("GetTeamAndPlayers: Error getting team\n" + err.Error())
    }
  }

  return all, team, nil
}
