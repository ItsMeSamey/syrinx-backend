package DB

import (
  "context"
  "log"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var DATABASE *mongo.Database

type Collection struct {
  coll *mongo.Collection
  context context.Context
}

// All the DB declarations
var (
  UserDB Collection
  QuestionDB Collection
) 

func InitDB(uri string) error {
  ctx := context.Background()
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
  if err != nil {
    return err
  }
  err = client.Ping(ctx, nil)
  if err != nil {
    return err
  }
  log.Println("Successfully Connected to MongoDB")

  DATABASE = client.Database("2024_ctf")
  UserDB = Collection{DATABASE.Collection("users"), context.TODO()}
  QuestionDB = Collection{DATABASE.Collection("questions"), context.TODO()}
  return nil
}


