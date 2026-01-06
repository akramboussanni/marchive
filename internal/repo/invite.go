package repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type InviteRepo struct {
	Columns
	db       *sqlx.DB
	userRepo *UserRepo
}

func NewInviteRepo(db *sqlx.DB, userRepo *UserRepo) *InviteRepo {
	repo := &InviteRepo{
		db:       db,
		userRepo: userRepo,
	}
	repo.Columns = ExtractColumns[model.Invite]()
	return repo
}

// CreateInvite creates a new invite for a user
func (r *InviteRepo) CreateInvite(ctx context.Context, inviterID int64) (*model.Invite, error) {
	// Generate unique token
	token, err := r.generateUniqueToken(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	// Create invite
	invite := &model.Invite{
		Token:     token,
		InviterID: inviterID,
		CreatedAt: now,
	}

	// Insert invite
	invite.ID = utils.GenerateSnowflakeID()
	query := fmt.Sprintf(`
		INSERT INTO invites (%s) 
		VALUES (%s)
	`, r.AllRaw, r.AllPrefixed)
	_, err = r.db.NamedExecContext(ctx, query, invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetInviteByToken gets an invite by its token
func (r *InviteRepo) GetInviteByToken(ctx context.Context, token string) (*model.Invite, error) {
	var invite model.Invite
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND revoked_at IS NULL`, r.AllRaw)
	err := r.db.GetContext(ctx, &invite, query, token)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

// UseInvite marks an invite as used and creates a new user
func (r *InviteRepo) UseInvite(ctx context.Context, token string, username string, passwordHash string) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get invite
	var invite model.Invite
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND revoked_at IS NULL AND used_at IS NULL`, r.AllRaw)
	err = tx.GetContext(ctx, &invite, query, token)
	if err != nil {
		return err
	}

	// Check if username already exists (using transaction connection)
	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username)
	if err != nil {
		return err
	}
	if exists {
		return ErrUsernameTaken
	}

	// Create user (using transaction connection)
	now := time.Now().Unix()
	user := &model.User{
		ID:             utils.GenerateSnowflakeID(),
		Username:       username,
		PasswordHash:   passwordHash,
		Role:           "user",
		CreatedAt:      now,
		JwtSessionID:   utils.GenerateSnowflakeID(),
		InviteTokens:   0,
		RequestCredits: 0,
	}

	userQuery := fmt.Sprintf(
		"INSERT INTO users (%s) VALUES (%s)",
		r.userRepo.AllRaw,
		r.userRepo.AllPrefixed,
	)
	_, err = tx.NamedExecContext(ctx, userQuery, user)
	if err != nil {
		return err
	}

	// Mark invite as used
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET invitee_username = $1, invitee_id = $2, used_at = $3 WHERE token = $4
	`, username, user.ID, now, token)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// RevokeInvite revokes an unused invite and returns the token to the user
func (r *InviteRepo) RevokeInvite(ctx context.Context, token string, inviterID int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get invite
	var invite model.Invite
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND inviter_id = $2 AND revoked_at IS NULL AND used_at IS NULL`, r.AllRaw)
	err = tx.GetContext(ctx, &invite, query, token, inviterID)
	if err != nil {
		return err
	}

	// Mark invite as revoked
	now := time.Now().Unix()
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET revoked_at = $1 WHERE token = $2
	`, now, token)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// GetUserInvites gets all invites for a user
func (r *InviteRepo) GetUserInvites(ctx context.Context, userID int64) ([]model.Invite, error) {
	var invites []model.Invite
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE inviter_id = $1 ORDER BY created_at DESC`, r.AllRaw)
	err := r.db.SelectContext(ctx, &invites, query, userID)
	return invites, err
}

// generateUniqueToken generates a unique invite token
func (r *InviteRepo) generateUniqueToken(ctx context.Context) (string, error) {
	for {
		// Generate 32-byte random token
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		token := hex.EncodeToString(bytes)

		// Check if token already exists
		var exists int
		err := r.db.GetContext(ctx, &exists, "SELECT 1 FROM invites WHERE token = $1", token)
		if err != nil {
			// Token doesn't exist, we can use it
			return token, nil
		}
		// Token exists, try again
	}
}

// Custom errors
var (
	ErrUsernameTaken = errors.New("username already taken")
)
