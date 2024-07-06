package DB

import (
  "context"
  "log"
  "errors"

  "go.mongodb.org/mongo-driver/bson"
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
  LobbyDB Collection
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
  LobbyDB = Collection{DATABASE.Collection("lobby"), context.TODO()}
  return nil
}

func (db *Collection) get(k, v string, out any) error {
  result := db.coll.FindOne(UserDB.context, bson.D{{k, v}})
  if result == nil {
    return errors.New("get: got a nil result")
  }
  if err := result.Err(); err != nil {
    return err
  }
  return result.Decode(out)
}

func (db *Collection) exists(k, v string) (bool, error) {
  result := UserDB.coll.FindOne(UserDB.context, bson.D{{k, v}})
  if result == nil {
    return false, errors.New("exists: got a nil result")
  }
  err := result.Err()
  if err == mongo.ErrNoDocuments {
    return false, nil
  }
  return err == nil, err
}

