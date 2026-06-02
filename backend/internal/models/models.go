package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        *string   `json:"email,omitempty"`
	Rating       int       `json:"rating"`
	GamesPlayed  int       `json:"gamesPlayed"`
	Karma        int       `json:"karma"`
	About        *string   `json:"about,omitempty"`
	Nationality  *string   `json:"nationality,omitempty"`
	Systems      []string  `json:"systems,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Tournament struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Game struct {
	ID             uuid.UUID  `json:"id"`
	Status         string     `json:"status"`
	ScoringType    string     `json:"scoringType"`
	TournamentID   *uuid.UUID `json:"tournamentId,omitempty"`
	NorthUserID    *uuid.UUID `json:"northUserId,omitempty"`
	EastUserID     *uuid.UUID `json:"eastUserId,omitempty"`
	SouthUserID    *uuid.UUID `json:"southUserId,omitempty"`
	WestUserID     *uuid.UUID `json:"westUserId,omitempty"`
	BoardNumber    int        `json:"boardNumber"`
	Seed           int64      `json:"seed"`
	DealPBN        string     `json:"dealPbn"`
	Auction        string     `json:"auction"`
	Contract       *string    `json:"contract,omitempty"`
	Declarer       *int       `json:"declarer,omitempty"`
	DeclarerTricks *int       `json:"declarerTricks,omitempty"`
	Score          *int       `json:"score,omitempty"`
	PlaySequence   string     `json:"playSequence"`
	CreatedAt      time.Time  `json:"createdAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
}
