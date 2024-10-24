package DB

import (
  "context"
  "errors"
  "log"
  "os"


  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "github.com/orcaman/concurrent-map/v2"
  utils "github.com/ItsMeSamey/go_utils"
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

  QUESTIONS = cmap.NewWithCustomShardingFunction[uint16, Question](func (v uint16) uint32 { return uint32(v) })
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

func tryHard(op func () error, maxTries byte) (err error) {
  var tries byte = 0
  sync:
  if err = utils.WithStack(op()); err != nil {
    if tries > maxTries { return }
    tries += 1;
    goto sync
  }
  return
}

func (db *Collection) replace(bsonM bson.M, entry any) error {
  result, err := db.Coll.ReplaceOne(db.Context, bsonM, entry);
  if err != nil {
    return utils.WithStack(err)
  }
  
  if result.MatchedCount == 0 {
    return utils.WithStack(errors.New("No document to be synced"))
  }

  return nil
}

func (db *Collection) replaceTryHard(bsonM bson.M, entry any, maxTries byte) error {
  return tryHard(func () error {
    return db.replace(bsonM, entry)
  }, maxTries)
}

/// Get the result of a db quarry in a `out` object, returns true if the object exists
/// NOTE: `out` must be a pointer or Programme will panic !
func (db *Collection) getExists(bsonM bson.M, out any) (bool, error) {
  result := db.Coll.FindOne(db.Context, bsonM)
  err := result.Err()

  if err == mongo.ErrNoDocuments {
    return false, nil
  } else if err != nil {
    return false, utils.WithStack(err)
  }

  if out != nil {
    if err := result.Decode(out); err != nil {
      return false, utils.WithStack(err)
    }
  }

  return true, nil
}

func (db *Collection) exists(bsonM bson.M) (bool, error) {
  return db.getExists(bsonM, nil)
}

func (db *Collection) get(bsonM bson.M, out any) error {
  exists, err := db.getExists(bsonM, out)
  if !exists {
    return utils.WithStack(errors.New("DB.get: document does not exist"))
  }
  return utils.WithStack(err)
}

