package DB

import (
  "time"
  "errors"
  "strings"

  "go.mongodb.org/mongo-driver/bson"
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
  err := QuestionDB.get("questionID", _id, &question)
  if err != nil {
    return nil, errors.New("QuestionFromID: DB.get error\n" + err.Error())
  }
  return &question, nil
}

func postQuestion(ques *Question) error {
  exists, err := QuestionDB.exists("question", ques.Question)
  if exists {
    return errors.New("postQuestion: Question already exists")
  }
  if err!=nil{
    return errors.New("postQuestion: Error in DB.exists" + err.Error())
  }
  _, err = QuestionDB.Coll.InsertOne(QuestionDB.Context, ques)
  if err != nil {
    return errors.New("postQuestion: Error occurred while adding question to database" + err.Error())
  }
  return nil
}

//check ans =ques id ,userid, answer 
func CheckAnswer(_id int16, userSessID SessID, givenAnswer string) (int, error) {
  question, err := QuestionFromID(_id)
  if err != nil {
    return 0, err
  }
  correct_answer := question.Answer
  isCorrect := strings.EqualFold(givenAnswer, correct_answer)
  if !isCorrect {
    return 0, nil
  }

  user, err := UserFromSessionID(userSessID)
  if err != nil {
    return 0,  err
  }
  if user == nil {
    return 0, errors.New("CheckAnswer: User not found")
  }

  team, err := UserByTeam(user.TeamID)
  if err != nil {
    return 0, err
  }
  if team.Solved[_id] >= 0 {
    return 0, errors.New("Question already solved")
  }
  team.Solved[_id] = time.Now().Unix()

  newPoints := len(team.Solved) * 100 * team.Level

  _, err = TeamDB.Coll.UpdateOne(
    TeamDB.Context,
    bson.M{"_id": team.TeamID},
    bson.M{
      "$set": bson.M{"Solved": team.Solved},
      "$inc": bson.M{"points": newPoints},
    },
    )
  if err != nil {
    return 0, err
  }

  // Fetch the updated team document to get the new points
  var updatedTeam struct {
    Points int `bson:"points"`
  }
  err = TeamDB.Coll.FindOne(
    TeamDB.Context,
    bson.M{"_id": team.TeamID},
    ).Decode(&updatedTeam)
  if err != nil {
    return 0, err
  }

  return updatedTeam.Points, nil
}

//point 
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

// return points - done
// remove genquestionid that will be set by us as 1 2 3 and so on - done
// findone to get -done
// teamfromid-done
// remove user -done
