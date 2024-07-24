package DB

import (
  "time"
  "errors"
  "strings"
  
  "go.mongodb.org/mongo-driver/bson"
)

/// Database sorted by TeamID
type Team struct {
  TeamID    TID             `bson:"teamID"`
  TeamName  string          `bson:"teamName"`
  Points    int             `bson:"points"`
  // Question id and time in unix milliseconds
  Solved    map[int16]int64 `bson:"solved"`
  // Question id and whether hint is used
  Hints     []int16         `bson:"hints"`
  Level     int             `bson:"level"`
  Exception bool            `bson:"exception"`
}

func createNewTeam(user *CreatableUser) error {
  _, err := TeamDB.Coll.InsertOne(TeamDB.Context, &Team{
    TeamID:    user.TeamID,
    TeamName:  user.TeamName,
    Points:    0,
    Solved:    make(map[int16]int64),
    Level:     0,
    Exception: false,
  })

  if err != nil {
    return errors.New("createTeam: Error while Team insertion" + err.Error())
  }

  return nil
}

func TeamByTeamID(teamID TID) (*Team, error) {
  var team Team
  if err := TeamDB.get(bson.M{"teamID": teamID}, &team); err != nil {
    return nil, errors.New("TeamByTeamID: DB.get failed\n"+err.Error())
  }
  return &team, nil
}

func (team *Team) IsSolved(ID int16) bool {
  _, ok := team.Solved[ID]
  return ok
}

func (team *Team) Sync(maxTries byte) error {
  if err := TeamDB.syncTryHard(bson.M{"teamID": team.TeamID}, team, maxTries); err != nil {
    return errors.New("Team.SyncTryHard: Error\n" + err.Error())
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

  if err := team.Sync(maxTries); err != nil {
    return "", errors.New(("Team.GetHint: sync error\n ") + err.Error())
  }

  return question.Hint, nil
}

/// Returns success(bool), error
func (team *Team) CheckAnswer(question *Question, Answer string, maxtries byte) (bool, error) {
  if team.IsSolved(question.ID) {
    return true, errors.New("Team.CheckAnswer: already solved")
  }
  if !strings.EqualFold(question.Answer, Answer) {
    return false, nil
  }

  team.Points += question.Points
  team.Solved[question.ID] = time.Now().UnixMilli()

  if err := team.Sync(maxtries); err != nil {
    return true, errors.New(("Team.CheckAnswer: sync error\n ") + err.Error())
  }

  return true, nil
}

