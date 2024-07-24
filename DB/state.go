package DB

import (
  "errors"
  "time"
  
  "go.mongodb.org/mongo-driver/bson"
)

type STATE struct {
  Level byte `bson:"level"`
}

var StateDB STATE


func stateSync() error {
  if err := ControlDB.get(bson.M{"type": "state"}, &StateDB); err != nil {
    return errors.New("stateSync: error in DB.get\n" + err.Error())
  }
  return nil
}

func initStateSynchronizer() {
  for {
    time.Sleep(5 * time.Second)
    stateSync()
  }
}

