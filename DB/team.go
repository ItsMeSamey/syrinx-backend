package DB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

/// Database sorted by TeamID
type Team struct {
  TeamID   TID `bson:"teamId"`
  TeamName  string  `bson:"teamName"`
  Points     int    `bson:"points"`
  Solved     map[primitive.ObjectID]int `bson:"solved"`
  Level      int    `bson:"level"`
}


func GetTeamNameByID(teamID TID) (string, error) {
  // Create a filter to find the team by its ID
  filter := bson.M{"teamID": teamID}

  // Create a result struct to hold the team name
  var result struct {
      TeamName string `bson:"TeamName"`
  }

  // Query the database
  err := TeamDB.Coll.FindOne(TeamDB.Context, filter).Decode(&result)
  if err != nil {
      if err == mongo.ErrNoDocuments {
          return "", errors.New("GetTeamNameByID: Team not found")
      }
      return "", err
  }

  return result.TeamName, nil
}




func createTeam(user *CreatableUser) error {

  newTeam := &Team{
    TeamID:   user.TeamID,
    TeamName: *user.TeamName,
    Points:   0,
    Solved:   make(map[primitive.ObjectID]int),
    Level:    1,
  }
  _, err := TeamDB.Coll.InsertOne(TeamDB.Context, newTeam)
  if err != nil {
      return  err
  }

  return nil
}

