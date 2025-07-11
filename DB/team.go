package DB

import (
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

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
func (team *Team) GetHint(question *Question, maxTries byte) string {
  for _, hint := range team.Hints {
    if hint == question.ID {
      return question.Hint
    }
  }
  team.Hints = append(team.Hints, question.ID)
  team.Points -= question.HintPoints

  return question.Hint
}

/// Returns success(bool), error
func (team *Team) CheckAnswer(question *Question, Answer string, maxtries byte) bool {
  if !strings.EqualFold(question.Answer, Answer) {
    return false
  }

  team.Solved[question.ID] = time.Now().UnixMilli()
  team.Points += question.Points

  return true
}

func (team *Team) Repoint() {
  points := 0

  for qid := range team.Solved {
    ques, err := GetQuestionFromIDTryHard(qid, 5)
    if err != nil {
      log.Println("Error for teamID: ", hex.EncodeToString((*team.TeamID)[:]), "\n", err)
    }
    points += ques.Points
  }

  for _, qid := range team.Hints {
    ques, err := GetQuestionFromIDTryHard(qid, 5)
    if err != nil {
      log.Println("Error for teamID: ", hex.EncodeToString((*team.TeamID)[:]), "\n", err)
    }
    points -= ques.HintPoints
  }

  err := team.Sync(5)
  if err != nil {
    log.Println("Error for teamID: ", hex.EncodeToString((*team.TeamID)[:]), "\n", err)
  }

  if team.Points != points {
    log.Println("teamID: ", hex.EncodeToString((*team.TeamID)[:]), ". Changing points from ", team.Points, " to ", points)
  }
  team.Points = points
}

