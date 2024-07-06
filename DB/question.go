package DB

import (
  "errors"

	"go.mongodb.org/mongo-driver/bson"
)

type Question struct {
  QuestionID string `bson:"questionID"`
  Question   string `bson:"question"`
  Points     int    `bson:"points"`
  Answer     string `bson:"answer"`
  Hint       string `bson:"hint"`
}

func GetQuestion(_id string) (*Question, error) {
	var question Question
	result := QuestionDB.coll.FindOne(UserDB.context, bson.D{{"_id", _id}})
	if result == nil {
		return nil, errors.New("UserFromSessionID: Token")
	}
	err := result.Decode(&question)
	if err != nil {
		return nil, err
	}
	return &question, err
}

