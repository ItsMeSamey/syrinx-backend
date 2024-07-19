package DB

import (
  "errors"
  "time"
  
  "go.mongodb.org/mongo-driver/bson"
)

/// Database sorted by TeamID
type Team struct {
  TeamID   TID            `bson:"teamID"`
  TeamName string         `bson:"teamName"`
  Points   int            `bson:"points"`
  // Question id and time in unix milliseconds
  Solved   map[int16]int64 `bson:"solved"`
  // Question id and whether hint is used
  // Hints     map[int16]bool `bson:"hints"`
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

func (team *Team) sync() error{
  result, err := TeamDB.Coll.ReplaceOne(TeamDB.Context, bson.M{"teamID": team.TeamID}, team);
  if err != nil {
    return errors.New("Error: Team.sync error\n" + err.Error())
  }
  
  if result.MatchedCount == 0 {
    return errors.New("Error: Team.sync failed\nmongod: No document found")
  }

  return nil
}

func (team *Team) SyncTryHard(maxTries byte) error {
  var tries byte = 0

  sync:
  if err := team.sync(); err != nil {
    if tries > maxTries {
      return errors.New("syncTryHard: Error in Team.Sync, Max Tries reached\n" + err.Error())
    }
    tries += 1;
    goto sync
  }

  return nil
}

/// Gives back the hint string
func (team *Team) GetHint(QID int16, maxTries byte) (string, error) {
  hint, points, err := GetHintTryHard(QID, maxTries)
  if err != nil {
    return "", errors.New("Error: Team.getHint\n" + err.Error())
  }

  // val, ok := team.Hints[QID]
  // if (!ok || val == false) {
  //   team.Hints[QID] = true
  //   team.Points -= points
  // }

  _ = points
  return hint, nil
}

/// Returns success(bool), error
func (team *Team) CheckAnswer(QID int16, Answer string, maxTries byte) (bool, error) {
  points, err := CheckAnswerTryHard(QID, Answer, maxTries)
  if err != nil {
    return false, errors.New("Error: Team.checkAnswer\n" + err.Error())
  }
  
  if (points == 0) {
    return false, nil
  }

  team.Points += points
  team.Solved[QID] = time.Now().UnixMilli()
  return true, nil
}

