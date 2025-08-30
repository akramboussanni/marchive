package repo

import (
	"context"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type FavoriteRepo struct {
	db *sqlx.DB
}

func NewFavoriteRepo(db *sqlx.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (fr *FavoriteRepo) AddFavorite(ctx context.Context, userID int64, bookHash string) error {
	now := time.Now().Unix()

	query := `
		INSERT INTO favorites (id, user_id, book_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`

	id := utils.GenerateSnowflakeID()
	_, err := fr.db.ExecContext(ctx, query, id, userID, bookHash, now)
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
	query := `
		SELECT f.id, f.user_id, f.book_hash, f.created_at
		FROM favorites f
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := fr.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []*model.Favorite
	for rows.Next() {
		fav := &model.Favorite{}
		err := rows.Scan(&fav.ID, &fav.UserID, &fav.BookHash, &fav.CreatedAt)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}

	return favorites, nil
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
