package DB

import (
  "time"
  "errors"
  
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
  }

  CallbackFunc func (prev, cur *STATE)
)

var (
  /// The var to store state
  State *STATE = &STATE{
    Level: 0,
  }
  Callbacks map[string]CallbackFunc = make(map[string]CallbackFunc)
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
    stateSync(bson.M{"type": "state", "changed": true})

    tries := byte(0)
    start:
    _, err := SyncDB.Coll.UpdateOne(SyncDB.Context, bson.M{"_id": State.id}, bson.D{{"$set", bson.M{"changed": false}}})
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
  go callCallbacks(prev, State)

  return nil
}

func callCallbacks(prev, cur *STATE) {
  for _, val := range Callbacks {
    val(prev, cur)
  }
}

