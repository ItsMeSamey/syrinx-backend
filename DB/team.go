package DB

import (

)

const (
	teamBucket = "teams"
)

type Team struct {
	TeamID    int
	Usernames []string
}


