package DB

import (
  "errors"
)

/// Database sorted by TeamID
type Team struct {
  TeamID   TID            `bson:"teamID"`
  TeamName string         `bson:"teamName"`
  Points   int            `bson:"points"`
  // Question id and time in unix milliseconds
  Solved   map[int16]int64 `bson:"solved"`
  // Question id and whether hint is used
  // Hint     map[int16]bool `bson:"hint"`
  Level    int            `bson:"level"`
}

func TeamNameByID(teamID TID) (string, error) {
  var result Team
  if err := TeamDB.get("teamID", teamID, &result); err != nil {
    return "", errors.New("TeamNameByID: DB.get failed\n"+err.Error())
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

