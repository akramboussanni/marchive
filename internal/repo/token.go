package repo

import (
	"context"
	"fmt"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/jmoiron/sqlx"
)

type TokenRepo struct {
	Columns
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) *TokenRepo {
	repo := &TokenRepo{db: db}
	repo.Columns = ExtractColumns[model.JwtBlacklist]()
	return repo
}

func (r *TokenRepo) RevokeToken(ctx context.Context, token model.JwtBlacklist) error {
	query := fmt.Sprintf(`
		INSERT INTO jwt_blacklist (%s)
		VALUES (%s)
		ON CONFLICT(jti) DO NOTHING
	`, r.AllRaw, r.AllPrefixed)
	_, err := r.db.NamedExecContext(ctx, query, token)
	return err
}

func (r *TokenRepo) IsTokenRevoked(jti string) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM jwt_blacklist WHERE jti = $1)
	`, jti)
	return exists, err
}

func (r *TokenRepo) CleanupTokens(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM jwt_blacklist WHERE expires_at < CURRENT_TIMESTAMP
	`)
	return err
}
