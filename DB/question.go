package DB

import (
	"errors"
	"math/rand"
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

func QuestionFromID(_id string) (*Question, error) {
	var question Question
	return &question, QuestionDB.get("_id", _id, &question)
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
