package handlers

import (
	"encoding/json"
	"game-matchmaking/backend/app"
	"game-matchmaking/backend/models"
	"net/http"
	"strconv"
	"time"
)

func GetGames(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		var games []models.Game
		query := app.DB.Preload("Maps")
		status := r.URL.Query().Get("status")
		if status != "" && status != "deleted" && status != "draft" {
			query = query.Where("status = ?", status)
		} else {
			query = query.Where("status NOT IN (?, ?)", "deleted", "draft")
		}
		var player models.Player
		if err := app.DB.First(&player, userID).Error; err == nil && player.Role == "creator" {
			query = query.Where("creator_id = ?", userID)
		}
		query.Find(&games)
		json.NewEncoder(w).Encode(games)
	}
}

func GetGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/games/"):]
		id, _ := strconv.Atoi(idStr)
		var game models.Game
		if err := app.DB.Preload("Maps").First(&game, id).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(game)
	}
}

func FormGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		if session.Values["user_id"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		idStr := r.URL.Path[len("/api/games/form/"):]
		id, _ := strconv.Atoi(idStr)
		var game models.Game
		if err := app.DB.First(&game, id).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		if game.Status != "draft" {
			http.Error(w, "Can only form draft games", http.StatusBadRequest)
			return
		}
		now := time.Now()
		game.Status = "formed"
		game.FormedAt = &now
		app.DB.Save(&game)
		json.NewEncoder(w).Encode(game)
	}
}

func CompleteGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var player models.Player
		if err := app.DB.First(&player, userID).Error; err != nil || player.Role != "moderator" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		idStr := r.URL.Path[len("/api/games/complete/"):]
		id, _ := strconv.Atoi(idStr)
		var game models.Game
		if err := app.DB.First(&game, id).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		if game.Status != "formed" {
			http.Error(w, "Can only complete formed games", http.StatusBadRequest)
			return
		}
		now := time.Now()
		game.Status = "completed"
		game.CompletedAt = &now
		game.ModeratorID = &player.ID
		app.DB.Save(&game)
		json.NewEncoder(w).Encode(game)
	}
}

func RejectGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var player models.Player
		if err := app.DB.First(&player, userID).Error; err != nil || player.Role != "moderator" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		idStr := r.URL.Path[len("/api/games/reject/"):]
		id, _ := strconv.Atoi(idStr)
		var game models.Game
		if err := app.DB.First(&game, id).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		if game.Status != "formed" {
			http.Error(w, "Can only reject formed games", http.StatusBadRequest)
			return
		}
		now := time.Now()
		game.Status = "rejected"
		game.CompletedAt = &now
		game.ModeratorID = &player.ID
		app.DB.Save(&game)
		json.NewEncoder(w).Encode(game)
	}
}

func ManageGames(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var games []models.Game
		app.DB.Preload("Maps").Find(&games)
		json.NewEncoder(w).Encode(games)
	}
}
