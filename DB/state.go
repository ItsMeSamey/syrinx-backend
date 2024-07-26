package DB

import (
  "errors"
  "time"
  
  "go.mongodb.org/mongo-driver/bson"
)

type (
  /// The struct that contains the state of the server
  STATE struct {
    // Private fields
    id      ObjID `bson:"_id,omitempty"`
    changed bool  `bson:"changed"`

    // Public fields
    Level int `bson:"level"`
    Keep  int `bson:"keep"`
    GameOn bool `bson:"gameOn"`
    SignOn bool `bson:"signOn"`
    TeamExceptions []any `bson:"teamExceptions"`
    Repoint bool `bson:"repointAll"`
  }

  CallbackFunc func (prev, cur *STATE)
)

var (
  /// The var to store state
  State *STATE = &STATE{
    Level: 0,
  }
  Callback CallbackFunc = nil
)

func InitSynchronizer() error {
  if err := stateSync(bson.M{"type": "state"}); err != nil {
    return errors.New("DB.Init: stateSync\n" + err.Error())
  }

  go startStateSynchronizer(5)
  return nil
}

func startStateSynchronizer(maxTries byte) {
  for {
    time.Sleep(2 * time.Second)
    err := stateSync(bson.M{"type": "state", "changed": true})
    if err != nil {
      continue
    }

    tries := byte(0)
    start:
    _, err = SyncDB.Coll.UpdateOne(SyncDB.Context, bson.M{"type": "state", "changed": true}, bson.D{{"$set", bson.M{"changed": false, "teamExceptions": nil}, }})
    if tries < maxTries && err != nil {
      tries  += 1
      goto start
    }
  }
}

func stateSync(bsonM bson.M) error {
  var NEW STATE
  prev := State

  if err := SyncDB.get(bsonM, &NEW); err != nil {
    return errors.New("stateSync: error in DB.get\n" + err.Error())
  }

  NEW.changed = false
  State = &NEW
  if Callback != nil {
    Callback(prev, State)
  }

  return nil
}

