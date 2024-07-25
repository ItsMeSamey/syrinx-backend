package DB

import (
  "os"
  "log"
  "sync"
  "errors"
  "context"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type (
  /// The Extendable Collection type
  Collection struct {
    Coll *mongo.Collection
    Context context.Context
  }

  /// Type declaration used for ID's
  TID *[3]byte
  SessID *[64]byte
  ObjID *primitive.ObjectID
)

var (
  /// The main Database
  DATABASE *mongo.Database

  /// All the DB declarations
  QuestionDB Collection
  UserDB     Collection
  TeamDB     Collection
  SyncDB     Collection

  /// Initialize env vars
  EMAIL_SENDER          = os.Getenv("EMAIL_SENDER")
  EMAIL_SENDER_PASSWORD = os.Getenv("EMAIL_SENDER_PASSWORD")

  QUESTIONS map[int16]Question = make(map[int16]Question)
  QUESTIONSMUTEX sync.RWMutex = sync.RWMutex{}
)

/// Initialize all Database's
/// Programme MUST panic if this function errors as this is unrecoverable
func Init() error {
  ctx := context.Background()
  uri := os.Getenv("MONGOURI")

  if uri == "" {
    return errors.New("MONGOURI not set")
  }
  if EMAIL_SENDER == "" {
    return errors.New("EMAIL_SENDER not set")
  }
  if EMAIL_SENDER_PASSWORD == "" {
    return errors.New("EMAIL_SENDER_PASSWORD not set")
  }

  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
  if err != nil {
    return errors.New("DB.Init: mongo.Connect\n" + err.Error())
  }
  
  if err = client.Ping(ctx, nil); err != nil {
    return errors.New("DB.Init: client.Ping\n" + err.Error())
  }

  log.Println("Successfully Connected to MongoDB")

  DATABASE = client.Database("2024_ctf")

  QuestionDB = Collection{DATABASE.Collection("questions"), ctx}
  UserDB     = Collection{DATABASE.Collection("users"),     ctx}
  TeamDB     = Collection{DATABASE.Collection("teams"),     ctx}
  SyncDB     = Collection{DATABASE.Collection("state"),   ctx}

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
  result := db.Coll.FindOne(db.Context, bsonM)
  err := result.Err()

  if err == mongo.ErrNoDocuments{
    return false, nil
  } else if err != nil {
    return false, errors.New("getExistsBson: DB.FindOne error\n" + err.Error())
  }

  if out != nil {
    if err := result.Decode(out); err != nil {
      return false, errors.New("getExistsBson: result.Decode error\n" + err.Error())
    }
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

