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

func QuestionFromID(ID int16) (*Question, error) {
  var question Question
  err := QuestionDB.get("questionID", ID, &question)
  if err != nil {
    return nil, errors.New("QuestionFromID: DB.get error\n" + err.Error())
  }
  return &question, nil
}

func GetQuestionFromIDTryHard(ID int16, maxTries byte) (*Question, error) {
  var tries byte = 0

  getQuestion:
  question, err := QuestionFromID(ID)
  if err != nil {
    if tries > maxTries {
      return nil, errors.New("GetQuestionFromIDTryHard: Error in DB.exists, Max Tries reached\n" + err.Error())
    }
    tries += 1;
    goto getQuestion
  }

  return question, nil
}

func postQuestion(ques *Question) error {
  exists, err := QuestionDB.exists("question", ques.Question)
  if exists {
    return errors.New("postQuestion: Question already exists")
  }
  if err!=nil{
    return errors.New("postQuestion: Error in DB.exists\n" + err.Error())
  }
  _, err = QuestionDB.Coll.InsertOne(QuestionDB.Context, ques)
  if err != nil {
    return errors.New("postQuestion: Error occurred while adding question to database\n" + err.Error())
  }
  return nil
}

//check ans = ques id ,userid, answer 
func CheckAnswerTryHard(ID int16, Answer string) (int, error) {
  question, err := GetQuestionFromIDTryHard(ID, 10)
  if err != nil {
    return 0, errors.New("CheckAnswerTryHard: Error while getting Question\n" + err.Error())
  }

  if strings.EqualFold(question.Answer, Answer) {
    return question.HintPoints, nil
  }
  
  return 0, nil
}

func GetHintTryHard(ID int16) (string, int, error) {
  question, err := GetQuestionFromIDTryHard(ID, 10)
  if err != nil {
    return  "", 0, errors.New("GetHintTryHard: Error while getting Question\n" + err.Error())
  }

  return question.Hint, question.HintPoints, nil
}

