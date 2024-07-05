package GdHandler

import (
	"ccs.ctf/DB"
)

func InitLobbies() error {
	DB.UserDB.ForEachInBucket()
	return nil
}
