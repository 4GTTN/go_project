package app

import (
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type App struct {
	DB    *gorm.DB
	Store *sessions.CookieStore
}

func NewApp(db *gorm.DB) *App {
	store := sessions.NewCookieStore([]byte("secret-key"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &App{
		DB:    db,
		Store: store,
	}
}
