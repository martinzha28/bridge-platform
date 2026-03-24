package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/game"
	"github.com/martinzha28/bridge-platform/backend/internal/models"
)

type GameRepository struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) *GameRepository {
	if db == nil {
		return nil
	}
	return &GameRepository{db: db}
}

func (r *GameRepository) SaveGame(ctx context.Context, g *models.Game) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO games (
			id, status, scoring_type, tournament_id,
			north_user_id, east_user_id, south_user_id, west_user_id,
			board_number, seed, deal_pbn, auction,
			contract, declarer, declarer_tricks, score,
			play_sequence, created_at, completed_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19
		)`,
		g.ID, g.Status, g.ScoringType, g.TournamentID,
		g.NorthUserID, g.EastUserID, g.SouthUserID, g.WestUserID,
		g.BoardNumber, g.Seed, g.DealPBN, g.Auction,
		g.Contract, g.Declarer, g.DeclarerTricks, g.Score,
		g.PlaySequence, g.CreatedAt, g.CompletedAt,
	)
	return err
}

func (r *GameRepository) GetGame(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	g := &models.Game{}
	err := r.db.QueryRow(ctx, `
		SELECT id, status, scoring_type, tournament_id,
			north_user_id, east_user_id, south_user_id, west_user_id,
			board_number, seed, deal_pbn, auction,
			contract, declarer, declarer_tricks, score,
			play_sequence, created_at, completed_at
		FROM games WHERE id = $1`, id,
	).Scan(
		&g.ID, &g.Status, &g.ScoringType, &g.TournamentID,
		&g.NorthUserID, &g.EastUserID, &g.SouthUserID, &g.WestUserID,
		&g.BoardNumber, &g.Seed, &g.DealPBN, &g.Auction,
		&g.Contract, &g.Declarer, &g.DeclarerTricks, &g.Score,
		&g.PlaySequence, &g.CreatedAt, &g.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *GameRepository) ListGamesByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.Game, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, status, scoring_type, tournament_id,
			north_user_id, east_user_id, south_user_id, west_user_id,
			board_number, seed, deal_pbn, auction,
			contract, declarer, declarer_tricks, score,
			play_sequence, created_at, completed_at
		FROM games
		WHERE north_user_id = $1 OR east_user_id = $1
		   OR south_user_id = $1 OR west_user_id = $1
		ORDER BY completed_at DESC NULLS LAST
		LIMIT $2`, userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(
			&g.ID, &g.Status, &g.ScoringType, &g.TournamentID,
			&g.NorthUserID, &g.EastUserID, &g.SouthUserID, &g.WestUserID,
			&g.BoardNumber, &g.Seed, &g.DealPBN, &g.Auction,
			&g.Contract, &g.Declarer, &g.DeclarerTricks, &g.Score,
			&g.PlaySequence, &g.CreatedAt, &g.CompletedAt,
		); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

// GameFromSession converts a completed game.Game into a models.Game for persistence.
func GameFromSession(g *game.Game) models.Game {
	now := time.Now()
	record := models.Game{
		ID:           uuid.New(),
		Status:       "completed",
		ScoringType:  "IMP",
		BoardNumber:  g.Board.Number,
		Seed:         g.Board.Seed,
		DealPBN:      g.Board.ToPBN(),
		Auction:      g.Auction.ToNotation(),
		PlaySequence: "",
		CreatedAt:    now,
		CompletedAt:  &now,
	}

	if g.Play != nil {
		record.PlaySequence = g.Play.ToNotation()
	}

	if g.Result != nil && !g.Result.PassedOut() {
		contract := g.Result.Contract.String()
		record.Contract = &contract

		declarer := int(g.Result.Contract.Declarer)
		record.Declarer = &declarer

		tricks := g.Result.TricksNS
		if g.Result.Contract.Declarer == game.East || g.Result.Contract.Declarer == game.West {
			tricks = g.Result.TricksEW
		}
		record.DeclarerTricks = &tricks
		record.Score = &g.Result.Score
	}

	return record
}
