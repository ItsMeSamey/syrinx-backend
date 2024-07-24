package DB

import (
  "errors"
  "time"
  
  "go.mongodb.org/mongo-driver/bson"
)

type (
  STATE struct {
    ID      ObjID `bson:"_id,omitempty"`
    Level   byte  `bson:"level"`
    Changed bool  `bson:"changed"`
  }
  CallbackFunc func (prev, cur *STATE)
)

var (
  State *STATE = nil
  Callbacks map[string]CallbackFunc = make(map[string]CallbackFunc)
)

func stateSync() error {
  var NEW STATE
  if err := SyncDB.get(bson.M{"type": "state", "changed": true}, &NEW); err != nil {
    return errors.New("stateSync: error in DB.get\n" + err.Error())
  }
  State = &NEW
  return nil
}

func callCallbacks(prev, cur *STATE) {
  for _, val := range Callbacks {
    val(prev, cur)
  }
}

func initStateSynchronizer(maxTries byte) {
  for {
    time.Sleep(5 * time.Second)

    prev := State
    stateSync()
    cur := State

    if prev != cur {
      cur.Changed = false
      go callCallbacks(prev, cur)
    }

    tries := byte(0)
    start:
    _, err := SyncDB.Coll.UpdateOne(SyncDB.Context, bson.M{"_id": cur.ID}, bson.A{"$set", bson.M{"changed": false}})
    if tries < maxTries && err != nil {
      tries  += 1
      goto start
    }
  }
}

