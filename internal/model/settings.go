package model

// AppSetting represents a key-value setting in the database
type AppSetting struct {
	Key       string `db:"key" json:"key"`
	Value     string `db:"value" json:"value"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
}

// Setting keys
const (
	SettingAnonymousAccessEnabled = "anonymous_access_enabled"
)
