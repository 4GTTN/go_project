// main.go
package main

import (
	"game-matchmaking/backend/app"
	"game-matchmaking/backend/handlers"
	"game-matchmaking/backend/models"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	dsn := "host=localhost user=postgres password=123456 dbname=game port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Map{}, &models.Player{}, &models.Game{}, &models.GameMap{})

	app := app.NewApp(db)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", enableCORS(handlers.ListHandler(app)))
	http.HandleFunc("/map/", enableCORS(handlers.DetailHandler(app)))
	http.HandleFunc("/delete/", enableCORS(handlers.RoleBasedMiddleware(handlers.DeleteHandler(app), []string{"admin"}, app)))
	http.HandleFunc("/api/maps", enableCORS(handlers.GetMaps(app)))
	http.HandleFunc("/api/maps/", enableCORS(handlers.GetMap(app)))
	http.HandleFunc("/api/maps/create", enableCORS(handlers.RoleBasedMiddleware(handlers.CreateMap(app), []string{"admin", "creator"}, app)))
	http.HandleFunc("/api/games", enableCORS(handlers.GetGames(app)))
	http.HandleFunc("/api/games/", enableCORS(handlers.GetGame(app)))
	http.HandleFunc("/api/games/add-map/", enableCORS(handlers.RoleBasedMiddleware(handlers.AddMapToGame(app), []string{"creator"}, app)))
	http.HandleFunc("/api/games/form/", enableCORS(handlers.RoleBasedMiddleware(handlers.FormGame(app), []string{"creator"}, app)))
	http.HandleFunc("/api/games/complete/", enableCORS(handlers.RoleBasedMiddleware(handlers.CompleteGame(app), []string{"moderator"}, app)))
	http.HandleFunc("/api/games/reject/", enableCORS(handlers.RoleBasedMiddleware(handlers.RejectGame(app), []string{"moderator"}, app)))
	http.HandleFunc("/api/games/remove-map/", enableCORS(handlers.RoleBasedMiddleware(handlers.RemoveMapFromGame(app), []string{"creator"}, app)))
	http.HandleFunc("/api/players/register", enableCORS(handlers.Register(app)))
	http.HandleFunc("/api/players/login", enableCORS(handlers.Login(app)))
	http.HandleFunc("/api/players/logout", enableCORS(handlers.Logout(app)))
	http.HandleFunc("/api/players/validate", enableCORS(handlers.ValidateSession(app)))
	http.HandleFunc("/api/manager/users", enableCORS(handlers.RoleBasedMiddleware(handlers.ManageUsers(app), []string{"admin"}, app)))
	http.HandleFunc("/api/manager/games", enableCORS(handlers.RoleBasedMiddleware(handlers.ManageGames(app), []string{"admin", "moderator"}, app)))
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
