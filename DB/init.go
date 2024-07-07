package DB

import (
  "context"
  "errors"
  "log"
  
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
  Coll *mongo.Collection
  Context context.Context
}

/// Type declaration used for ID's
type TID *[3]byte
type SessID *[64]byte
type ObjID *primitive.ObjectID

/// The main Database
var DATABASE *mongo.Database

/// All the DB declarations
var (
  UserDB Collection
  QuestionDB Collection
  LobbyDB Collection
) 

/// Initialize all Database's
/// Programme MUST panic if this function errors as this is unrecoverable
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

/// Get the result of a db quarry in a `out` object
/// NOTE: `out` must be a pointer or Programme will panic !
func (db *Collection) get(k string, v any, out any) error {
  result := db.Coll.FindOne(UserDB.Context, bson.D{{k, v}})
  if result == nil {
    return errors.New("get: got a nil result")
  }
  if err := result.Err(); err != nil {
    return err
  }
  return result.Decode(out)
}

/// Check if a entry exists in a Collection
func (db *Collection) exists(k string, v any) (bool, error) {
  result := UserDB.Coll.FindOne(UserDB.Context, bson.D{{k, v}})
  if result == nil {
    return false, errors.New("exists: got a nil result")
  }
  err := result.Err()
  if err == mongo.ErrNoDocuments {
    return false, nil
  }
  return err == nil, err
}

