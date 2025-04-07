package models

import (
	"time"

	"github.com/google/uuid"
)

// NOTE: fetch the MatchmakingRating from the user service
type Player struct {
	ID                string
	Username          string  `json:"username"`
	MatchmakingRating float64 `json:"match_making_rating"`
	QueueID           string  `json:"queue_id"`
}

func (p *Player) NewPlayer(username, queue_id string, mmr float64) *Player {
	return &Player{
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

type Match struct {
	ID        string
	Players   []Player `json:"players"`
	CreatedAt uint64   `json:"created_at"`
}

func (m Match) NewMatch(involved_players []Player) Match {
	return Match{
		ID:        uuid.NewString(),
		Players:   involved_players,
		CreatedAt: uint64(time.Now().Unix()),
	}
}
