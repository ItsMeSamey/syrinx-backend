package DB

import (
  "errors"
  "time"
  "strings"
  
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
  Hints     []int16 `bson:"hints"`
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
      return errors.New("Team.SyncTryHard: Error in Team.Sync, Max Tries reached\n" + err.Error())
    }
    tries += 1;
    goto sync
  }

  return nil
}

/// Gives back the hint string
func (team *Team) GetHint(question *Question, maxTries byte) (string, error) {
  for _, hint := range team.Hints {
    if hint == question.ID {
      return question.Hint, nil
    }
  }
  team.Hints = append(team.Hints, question.ID)
  team.Points -= question.HintPoints

  if err := team.SyncTryHard(maxTries); err != nil {
    return "", errors.New(("Team.GetHint: sync error\n ") + err.Error())
  }

  return question.Hint, nil
}

/// Returns success(bool), error
func (team *Team) CheckAnswer(question *Question, Answer string, maxtries byte) (bool, error) {

  if _, ok := team.Solved[question.ID]; ok {
    return true, errors.New("Team.CheckAnswer: already solved")
  }

  if !strings.EqualFold(question.Answer, Answer) {
    return false, nil
  }

  team.Points += question.Points
  team.Solved[question.ID] = time.Now().UnixMilli()

  if err := team.SyncTryHard(maxtries); err != nil {
    return true, errors.New(("Team.CheckAnswer: sync error\n ") + err.Error())
  }

  return true, nil
}

