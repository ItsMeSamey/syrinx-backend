package DB

import (
  "errors"
	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
)

/// Database sorted by TeamID
type Team struct {
  TeamID   TID            `bson:"teamID"`
  TeamName string         `bson:"teamName"`
  Points   int            `bson:"points"`
  // Question id and time in unix milliseconds
  Solved   map[int16]int64 `bson:"solved"`
  // Question id and whether hint is used
  Hint     map[int16]bool `bson:"hint"`
  Level    int            `bson:"level"`
}
func UserByTeam(teamID TID) (*Team, error) {
  var team Team
  err := TeamDB.get("teamID", teamID, &team)
  if err != nil {
      if err == mongo.ErrNoDocuments {
          return nil, errors.New("no team found with the given teamID")
      }
      return nil, err
  }
  return &team, nil
}


func TeamNameByID(teamID TID) (string, error) {
  var result Team
  if err := TeamDB.get("teamID", teamID, &result); err != nil {
    return "", errors.New("getTeamNameByID: DB.get failed\n"+err.Error())
  }
  return result.TeamName, nil
}

func createNewTeam(user *CreatableUser) error {
  _, err := TeamDB.Coll.InsertOne(TeamDB.Context, &Team{
    TeamID:   user.TeamID,
    TeamName: *user.TeamName,
    Points:   0,
    Solved:   make(map[int16]int64),
    Level:    0,
  })

  if err != nil {
    return errors.New("createTeam: Error while Team insertion" + err.Error())
  }

  return nil
}

