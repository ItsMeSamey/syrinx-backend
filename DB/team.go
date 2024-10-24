package DB

import (
  "encoding/hex"
  "log"
  "strings"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  utils "github.com/ItsMeSamey/go_utils"
)

/// Database sorted by TeamID
type Team struct {
  TeamID    TID             `bson:"teamID"`
  TeamName  string          `bson:"teamName"`
  Points    int             `bson:"points"`
  // Question id and time in unix milliseconds
  Solved    map[uint16]int64 `bson:"solved"`
  // Question id and whether hint is used
  Hints     []uint16        `bson:"hints"`
  Level     int             `bson:"level"`
  Exception bool            `bson:"exception"`
}

func createNewTeam(user *CreatableUser) error {
  _, err := TeamDB.Coll.InsertOne(TeamDB.Context, &Team{
    TeamID:    user.TeamID,
    TeamName:  user.TeamName,
    Points:    0,
    Solved:    make(map[uint16]int64),
    Level:     0,
    Exception: false,
  })

  return utils.WithStack(err)
}

func TeamByTeamID(teamID TID) (team *Team, err error) {
  err = utils.WithStack(TeamDB.get(bson.M{"teamID": teamID}, team))
  return
}

func (team *Team) IsSolved(ID uint16) (ok bool) {
  _, ok = team.Solved[ID]
  return
}

func (team *Team) Sync(maxTries byte) error {
  return TeamDB.replaceTryHard(bson.M{"teamID": team.TeamID}, team, maxTries)
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
    ques, err := QuestionFromID(qid, 5)
    if err != nil {
      log.Println("Error for teamID: ", hex.EncodeToString((*team.TeamID)[:]), "\n", err)
    }
    points += ques.Points
  }

  for _, qid := range team.Hints {
    ques, err := QuestionFromID(qid, 5)
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

