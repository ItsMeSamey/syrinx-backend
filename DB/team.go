package DB

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
  ID         ObjID  `bson:"_id,omitempty"`
  TeamID   string `bson:"question"`
  TeamName  string  `bson:"teamName"`
  Points     int    `bson:"points"`
  Solved     map[primitive.ObjectID]int `bson:"solved"`
  Level      int    `bson:"level"`
}

