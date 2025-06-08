// app.go
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
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Adjust based on your needs
	}
	return &App{
		DB:    db,
		Store: store,
	}
}
