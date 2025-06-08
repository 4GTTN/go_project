package main

import (
	"backend/app"
	"backend/handlers"
	"backend/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
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
func seedDatabase(db *gorm.DB) error {
	var mapCount int64
	db.Model(&models.Map{}).Count(&mapCount)
	if mapCount == 0 {
		maps := []models.Map{
			{
				Name:        "Desert Arena",
				Price:       10.99,
				Description: "A sandy battleground with ancient ruins",
				Location:    "Desert",
				Image:       "https://i.pinimg.com/originals/d4/11/96/d411960a2ae0a67117db7016f8142734.jpg",
				IsDeleted:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Urban City",
				Price:       15.99,
				Description: "A modern cityscape for intense battles",
				Location:    "City",
				Image:       "https://i.pinimg.com/736x/c3/be/77/c3be771408ef77b215737db7ca12027c.jpg",
				IsDeleted:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "snowy",
				Price:       12.99,
				Description: "snow-capped mountains with caves",
				Location:    "snow",
				Image:       "https://mir-s3-cdn-cf.behance.net/project_modules/1400_opt_1/4a010733810574.56b91309db01c.jpg",
				IsDeleted:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "sky sity",
				Price:       12.99,
				Description: "Cloud City from the Star Wars universe",
				Location:    "City",
				Image:       "https://i.ytimg.com/vi/YZEWHIwmmIM/maxresdefault.jpg",
				IsDeleted:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "china-town",
				Price:       12.99,
				Description: "A small Chinese town like in the anime",
				Location:    "City",
				Image:       "https://i.pinimg.com/originals/e4/7c/e0/e47ce0caf966297effbac58d34848e29.jpg",
				IsDeleted:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		if err := db.Create(&maps).Error; err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Формируем строку подключения из переменных окружения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Map{}, &models.Player{}, &models.Game{}, &models.GameMap{})

	if err := seedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
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
