package internal

import "github.com/google/uuid"

type Player struct {
	ID                string
	Username          string  `json:"username"`
	MatchmakingRating float64 `json:"match_making_rating"`
	QueueID           string  `json:"queue_id"`
}

func (p Player) NewPlayer(username, queue_id string, mmr float64) Player {
	return Player{
		ID:                uuid.NewString(),
		Username:          username,
		MatchmakingRating: mmr,
		QueueID:           queue_id,
	}
}

type AbstractPlayer struct {
	ID    string
	Score float64
}
