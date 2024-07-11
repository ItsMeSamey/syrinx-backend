package DB

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

/// Database sorted by TeamID
type Team struct {
  TeamID   TID `bson:"question"`
  TeamName  string  `bson:"teamName"`
  Points     int    `bson:"points"`
  Solved     map[primitive.ObjectID]int `bson:"solved"`
  Level      int    `bson:"level"`
}


func getTeamNameByTID(id TID) (string, error) {
  return "", nil
}

