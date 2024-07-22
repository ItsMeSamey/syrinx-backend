package DB

import (
  "os"
  "log"
  "errors"
  "context"
  
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson/primitive"
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
  QuestionDB Collection
  UserDB Collection
  TeamDB Collection
  LobbyDB Collection
) 

/// Initialize env vars
var (
  EMAIL_SENDER = os.Getenv("EMAIL_SENDER")
  EMAIL_SENDER_PASSWORD = os.Getenv("EMAIL_SENDER_PASSWORD")
)

/// Initialize all Database's
/// Programme MUST panic if this function errors as this is unrecoverable
func InitDB(uri string) error {
  if EMAIL_SENDER == "" {
    return errors.New("Email sender does not exist")
  }
  if EMAIL_SENDER_PASSWORD == "" {
    return errors.New("Email password not set")
  }
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
  QuestionDB = Collection{DATABASE.Collection("questions"), context.TODO()}
  UserDB = Collection{DATABASE.Collection("users"), context.TODO()}
  TeamDB = Collection{DATABASE.Collection("teams"), context.TODO()}
  LobbyDB = Collection{DATABASE.Collection("lobby"), context.TODO()}
  return nil
}

func (db *Collection) sync(bsonM bson.M, entry any) error {
  result, err := db.Coll.ReplaceOne(db.Context, bsonM, entry);
  if err != nil {
    return errors.New("Error: DB.syncBson error\n" + err.Error())
  }
  
  if result.MatchedCount == 0 {
    return errors.New("Error: DB.syncBson failed, No document synced")
  }

  return nil
}

func (db *Collection) syncTryHard(bsonM bson.M, entry any, maxTries byte) error {
  var tries byte = 0

  sync:
  if err := db.sync(bsonM, entry); err != nil {
    if tries > maxTries {
      return errors.New("DB.syncBsonTryHard: Error in DB.syncBson, Max Tries reached\n" + err.Error())
    }
    tries += 1;
    goto sync
  }

  return nil
}

/// Get the result of a db quarry in a `out` object, returns true if the object exists
/// NOTE: `out` must be a pointer or Programme will panic !
func (db *Collection) getExists(bsonM bson.M, out any) (bool, error) {
  result := db.Coll.FindOne(LobbyDB.Context, bsonM)
  err := result.Err()

  if err == mongo.ErrNoDocuments{
    return false, nil
  } else if err != nil {
    return false, errors.New("getExistsBson: DB.FindOne error\n" + err.Error())
  }

  if out == nil {
    return true, nil
  }

  if err := result.Decode(out); err != nil {
    return false, errors.New("getExistsBson: result.Decode error\n" + err.Error())
  }

  return true, nil
}

func (db *Collection) get(bsonM bson.M, out any) error {
  exists, err := db.getExists(bsonM, out)
  if !exists {
    return errors.New("DB.get: document does not exist")
  }
  return err
}

func (db *Collection) exists(bsonM bson.M) (bool, error) {
  return db.getExists(bsonM, nil)
}

