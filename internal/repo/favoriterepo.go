package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type FavoriteRepo struct {
	Columns
	db *sqlx.DB
}

func NewFavoriteRepo(db *sqlx.DB) *FavoriteRepo {
	repo := &FavoriteRepo{db: db}
	repo.Columns = ExtractColumns[model.Favorite]()
	return repo
}

func (fr *FavoriteRepo) AddFavorite(ctx context.Context, userID int64, bookHash string) error {
	now := time.Now().Unix()

	favorite := &model.Favorite{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		BookHash:  bookHash,
		CreatedAt: now,
	}

	query := fmt.Sprintf(`
		INSERT INTO favorites (%s)
		VALUES (%s)
	`, fr.AllRaw, fr.AllPrefixed)

	_, err := fr.db.NamedExecContext(ctx, query, favorite)
	return err
}

func (fr *FavoriteRepo) RemoveFavorite(ctx context.Context, userID int64, bookHash string) error {
	query := `
		DELETE FROM favorites 
		WHERE user_id = $1 AND book_hash = $2
	`

	_, err := fr.db.ExecContext(ctx, query, userID, bookHash)
	return err
}

func (fr *FavoriteRepo) IsFavorited(ctx context.Context, userID int64, bookHash string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM favorites 
		WHERE user_id = $1 AND book_hash = $2
	`

	var count int
	err := fr.db.QueryRowContext(ctx, query, userID, bookHash).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (fr *FavoriteRepo) GetUserFavorites(ctx context.Context, userID int64, limit, offset int) ([]*model.Favorite, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM favorites f
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`, fr.AllRaw)

	var favorites []*model.Favorite
	err := fr.db.SelectContext(ctx, &favorites, query, userID, limit, offset)
	return favorites, err
}

func (fr *FavoriteRepo) CountUserFavorites(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM favorites 
		WHERE user_id = $1
	`

	var count int
	err := fr.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
