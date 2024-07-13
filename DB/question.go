package DB

import (
	"errors"
	"math/rand"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"time"

)


type Question struct {
	ID         int16  `bson:"questionID"`
	Question   string `bson:"question"`
	Answer     string `bson:"answer"`
	Points     int    `bson:"points"`
	Hint       string `bson:"hint"` //need to change it to array of hints
	HintPoints int    `bson:"hintpoints"`
	Level      int    `bson:"level"`
}

func QuestionFromID(_id int16) (*Question, error) {
	var question Question
	err:=QuestionDB.get("questionID", _id, &question)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func genQuestionID() (int16, error) {
	times := 0
start:
	ID := int16(rand.Intn(32767))
	exists, err := QuestionDB.exists("questionID", ID)
	if exists {
		if times > 1024 {
			return 0, errors.New("genQuestionID: Lucky Error!!")
		}
		times += 1
		goto start
	}
	return ID, err
}

func postQuestion(ques *Question) (string, error) {
	ques.ID = -1
	exists, err := QuestionDB.exists("question", ques.Question)
	if exists {
		return "Question already exists", nil
	}
	ques.ID, _ = genQuestionID()
	if ques.ID == -1 {
		return "Error in generating question ID", nil
	}
	_, err = QuestionDB.Coll.InsertOne(QuestionDB.Context, &Question{
		ID:         ques.ID,
		Question:   ques.Question,
		Answer:     ques.Answer,
		Points:     ques.Points,
		Hint:       ques.Hint,
		HintPoints: ques.HintPoints,
		Level:      ques.Level,
	})
	return "ok", err
}


//check ans =ques id ,userid, answer 
func CheckAnswer(_id int16, userSessID SessID, givenAnswer string) (bool, error){
  question, err := QuestionFromID(_id)
  if err != nil {
      return false, err
  }
  correct_answer := question.Answer
  isCorrect := strings.EqualFold(givenAnswer, correct_answer)
  if isCorrect {
	user,err:=UserFromSessionID(userSessID)
    if user == nil {
		return false, errors.New("CheckAnswer: User not found")
	}
	team,err:=UserByTeam(user.TeamID)
	if err!=nil{
		return false,err
	}
	if team.Solved[_id]>=0{
		return false,errors.New("Question already solved")
	}
	team.Solved[_id]=time.Now().Unix()
	_, err = TeamDB.Coll.UpdateOne(
    TeamDB.Context,
    bson.M{"_id": team.TeamID},
    bson.M{"$set": bson.M{"Solved": team.Solved}},
    )
	if err != nil {
		return false, err
	}
	err=UpdateTeamPoints(team.TeamID,len(team.Solved)*100*team.Level)
	if err!=nil{
		return false,err
	}
	return true,nil
	}


}

//get just ques
func GetQuestionString(_id int16)(string,error){
	question, err := QuestionFromID(_id)
  	if err != nil {
    	  return "", err
  	}
	return question.Question,nil
}


//point 
func UpdateTeamPoints(teamID TID, newPoints int) error {
    _, err := TeamDB.Coll.UpdateOne(
        TeamDB.Context,
        bson.M{"_id": teamID},
        bson.M{"$set": bson.M{"points": newPoints}},
    )
    return err
}

//get hint
func GetHint(quesid int16,userSessID SessID)(string,error){
	question, err := QuestionFromID(quesid)
  	if err != nil {
		  return "", err
  	}
	user,err:=UserFromSessionID(userSessID)
	if user == nil {
		return "", errors.New("GetHint: User not found")
	}
	team,err:=UserByTeam(user.TeamID)
	if err!=nil{
		return "",err
	}
	if team.Hint[quesid]==true{
		return question.Hint,nil
	}
	if team.Points<30{
		return "",errors.New("Teri Aukat Nahi Hai!!!")
	}
 	_, err = TeamDB.Coll.UpdateOne(
    TeamDB.Context,
    bson.M{"_id": team.TeamID},
    bson.M{"$set": bson.M{"points": team.Points-30}},
)
 	if err != nil {
		return "", err
	}

	team.Hint[quesid]=true
	_, err = TeamDB.Coll.UpdateOne(
	  TeamDB.Context,
	  bson.M{"_id": team.TeamID},
	  bson.M{"$set": bson.M{"Hint": team.Hint}},
	)
	if err != nil {
		return "", err
	}

	return question.Hint,nil
}

// retun points
// findone to get
// teamfromid-done
// remove user -done
