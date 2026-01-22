package repo

import (
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repos struct {
	User              *UserRepo
	Token             *TokenRepo
	Lockout           *LockoutRepo
	DownloadRequest   *DownloadRequestRepo
	AnonymousDownload *AnonymousDownloadRepo
	Book              *BookRepo
	DownloadJob       *DownloadJobRepo
	SearchCache       *SearchCacheRepo
	Favorite          *FavoriteRepo
	RequestCredits    *RequestCreditsRepo
	Invite            *InviteRepo
	Settings          *SettingsRepo
}

type Columns struct {
	allColumns   []string
	AllRaw       string
	AllPrefixed  string
	safeColumns  []string
	SafeRaw      string
	SafePrefixed string
}

func NewRepos(db *sqlx.DB) *Repos {
	userRepo := NewUserRepo(db)

	return &Repos{
		User:              userRepo,
		Token:             NewTokenRepo(db),
		Lockout:           NewLockoutRepo(db),
		DownloadRequest:   NewDownloadRequestRepo(db),
		AnonymousDownload: NewAnonymousDownloadRepo(db),
		Book:              NewBookRepo(db),
		DownloadJob:       NewDownloadJobRepo(db),
		SearchCache:       NewSearchCacheRepo(db),
		Favorite:          NewFavoriteRepo(db),
		RequestCredits:    NewRequestCreditsRepo(db),
		Invite:            NewInviteRepo(db, userRepo),
		Settings:          NewSettingsRepo(db),
	}
}

func ExtractColumns[T any]() Columns {
	var allCols, safeCols []string

	t := reflect.TypeOf((*T)(nil)).Elem()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag, ok := field.Tag.Lookup("db")
		if !ok || dbTag == "-" {
			dbTag = strings.ToLower(field.Name)
		}
		allCols = append(allCols, dbTag)

		if safeTag, ok := field.Tag.Lookup("safe"); ok && safeTag == "true" {
			safeCols = append(safeCols, dbTag)
		}
	}

	allInsert := strings.Join(allCols, ", ")
	allSelect := ":" + strings.Join(allCols, ", :")

	safeInsert := strings.Join(safeCols, ", ")
	safeSelect := ":" + strings.Join(safeCols, ", :")

	return Columns{
		allColumns:   allCols,
		AllRaw:       allInsert,
		AllPrefixed:  allSelect,
		safeColumns:  safeCols,
		SafeRaw:      safeInsert,
		SafePrefixed: safeSelect,
	}
}
