package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type RedemptionCodeRepo struct {
	CodeColumns Columns
	LogColumns  Columns
	db          *sqlx.DB
}

func NewRedemptionCodeRepo(db *sqlx.DB) *RedemptionCodeRepo {
	repo := &RedemptionCodeRepo{db: db}
	repo.CodeColumns = ExtractColumns[model.RedemptionCode]()
	repo.LogColumns = ExtractColumns[model.RedemptionLog]()
	return repo
}

// CreateRedemptionCode creates a new redemption code
func (r *RedemptionCodeRepo) CreateRedemptionCode(ctx context.Context, req *model.CreateRedemptionCodeRequest, createdBy int64) (*model.RedemptionCode, error) {
	now := time.Now().Unix()
	
	code := &model.RedemptionCode{
		ID:             utils.GenerateSnowflakeID(),
		Code:           req.Code,
		Description:    req.Description,
		InviteTokens:   req.InviteTokens,
		RequestCredits: req.RequestCredits,
		MaxUses:        req.MaxUses,
		CurrentUses:    0,
		ExpiresAt:      req.ExpiresAt,
		RevokedAt:      nil,
		CreatedBy:      createdBy,
		CreatedAt:      now,
	}

	query := fmt.Sprintf(
		"INSERT INTO redemption_codes (%s) VALUES (%s)",
		r.CodeColumns.AllRaw,
		r.CodeColumns.AllPrefixed,
	)

	_, err := r.db.NamedExecContext(ctx, query, code)
	if err != nil {
		return nil, fmt.Errorf("failed to create redemption code: %w", err)
	}

	return code, nil
}

// GetRedemptionCodeByCode retrieves a redemption code by its code string
func (r *RedemptionCodeRepo) GetRedemptionCodeByCode(ctx context.Context, codeStr string) (*model.RedemptionCode, error) {
	query := fmt.Sprintf("SELECT %s FROM redemption_codes WHERE code = $1", r.CodeColumns.AllRaw)

	var code model.RedemptionCode
	err := r.db.GetContext(ctx, &code, query, codeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get redemption code: %w", err)
	}

	return &code, nil
}

// IsCodeRedeemedByUser checks if a user has already redeemed a specific code
func (r *RedemptionCodeRepo) IsCodeRedeemedByUser(ctx context.Context, codeID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM redemption_log WHERE code_id = $1 AND user_id = $2`
	
	var count int
	err := r.db.GetContext(ctx, &count, query, codeID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check if code was redeemed: %w", err)
	}
	
	return count > 0, nil
}

// RedeemCode marks a code as redeemed by a user and logs the redemption
func (r *RedemptionCodeRepo) RedeemCode(ctx context.Context, codeID, userID int64, inviteTokens, requestCredits int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Increment current uses count
	updateQuery := `UPDATE redemption_codes SET current_uses = current_uses + 1 WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, codeID)
	if err != nil {
		return fmt.Errorf("failed to update code usage count: %w", err)
	}

	// Log the redemption
	logQuery := fmt.Sprintf(
		"INSERT INTO redemption_log (%s) VALUES (%s)",
		r.LogColumns.AllRaw,
		r.LogColumns.AllPrefixed,
	)
	
	logEntry := &model.RedemptionLog{
		ID:                    utils.GenerateSnowflakeID(),
		CodeID:                codeID,
		UserID:                userID,
		RedeemedAt:            time.Now().Unix(),
		InviteTokensGranted:   inviteTokens,
		RequestCreditsGranted: requestCredits,
	}
	
	_, err = tx.NamedExecContext(ctx, logQuery, logEntry)
	if err != nil {
		return fmt.Errorf("failed to log redemption: %w", err)
	}

	// Update user's tokens and credits
	if inviteTokens > 0 {
		updateUserQuery := `UPDATE users SET invite_tokens = invite_tokens + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateUserQuery, inviteTokens, userID)
		if err != nil {
			return fmt.Errorf("failed to update user invite tokens: %w", err)
		}
	}

	if requestCredits > 0 {
		updateUserQuery := `UPDATE users SET request_credits = request_credits + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateUserQuery, requestCredits, userID)
		if err != nil {
			return fmt.Errorf("failed to update user request credits: %w", err)
		}
	}

	return tx.Commit()
}

// ListRedemptionCodes retrieves all redemption codes with optional pagination
func (r *RedemptionCodeRepo) ListRedemptionCodes(ctx context.Context, limit, offset int) ([]*model.RedemptionCode, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM redemption_codes 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, r.CodeColumns.AllRaw)

	var codes []*model.RedemptionCode
	err := r.db.SelectContext(ctx, &codes, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list redemption codes: %w", err)
	}

	return codes, nil
}

// CountRedemptionCodes returns the total number of redemption codes
func (r *RedemptionCodeRepo) CountRedemptionCodes(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM redemption_codes`
	
	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count redemption codes: %w", err)
	}
	
	return count, nil
}

// RevokeRedemptionCode marks a code as revoked
func (r *RedemptionCodeRepo) RevokeRedemptionCode(ctx context.Context, codeID int64) error {
	now := time.Now().Unix()
	
	query := `UPDATE redemption_codes SET revoked_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, now, codeID)
	if err != nil {
		return fmt.Errorf("failed to revoke redemption code: %w", err)
	}
	
	return nil
}

// DeleteRedemptionCode permanently deletes a redemption code
func (r *RedemptionCodeRepo) DeleteRedemptionCode(ctx context.Context, codeID int64) error {
	query := `DELETE FROM redemption_codes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, codeID)
	if err != nil {
		return fmt.Errorf("failed to delete redemption code: %w", err)
	}
	
	return nil
}
