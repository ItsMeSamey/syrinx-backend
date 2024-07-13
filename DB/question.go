package DB

import (
  "errors"
  "strings"
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
func CheckAnswer(ID int16, Answer string) (int, error) {
  const maxTries byte = 10
  var tries byte = 0

  getQuestion:
  question, err := QuestionFromID(ID)
  if err != nil {
    if tries > maxTries {
      return 0, errors.New("QuestionFromID: Error in DB.exists" + err.Error())
    }
    tries += 1;
    goto getQuestion
  }

  if strings.EqualFold(question.Answer, Answer) {
    return question.HintPoints, nil
  }
  
  return 0, nil
}

func GetHint(ID int16) (int, string, error) {
  const maxTries byte = 10
  var tries byte = 0

  getQuestion:
  question, err := QuestionFromID(ID)
  if err != nil {
    if tries > maxTries {
      return 0, "", errors.New("QuestionFromID: Error in DB.exists" + err.Error())
    }
    tries += 1;
    goto getQuestion
  }

  return question.HintPoints, question.Hint, nil
}

