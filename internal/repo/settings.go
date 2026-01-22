package repo

import (
	"context"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/jmoiron/sqlx"
)

type SettingsRepo struct {
	db *sqlx.DB
}

func NewSettingsRepo(db *sqlx.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

// GetSetting retrieves a setting by key
func (r *SettingsRepo) GetSetting(ctx context.Context, key string) (*model.AppSetting, error) {
	var setting model.AppSetting
	query := `SELECT key, value, updated_at FROM app_settings WHERE key = $1`
	err := r.db.GetContext(ctx, &setting, query, key)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SetSetting updates or inserts a setting
func (r *SettingsRepo) SetSetting(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO app_settings (key, value, updated_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = $3
	`
	_, err := r.db.ExecContext(ctx, query, key, value, time.Now().Unix())
	return err
}

// GetAllSettings retrieves all settings
func (r *SettingsRepo) GetAllSettings(ctx context.Context) ([]model.AppSetting, error) {
	var settings []model.AppSetting
	query := `SELECT key, value, updated_at FROM app_settings ORDER BY key`
	err := r.db.SelectContext(ctx, &settings, query)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// IsAnonymousAccessEnabled checks if anonymous access is enabled
func (r *SettingsRepo) IsAnonymousAccessEnabled(ctx context.Context) bool {
	setting, err := r.GetSetting(ctx, model.SettingAnonymousAccessEnabled)
	if err != nil {
		return false // Default to disabled if setting not found
	}
	return setting.Value == "true"
}
