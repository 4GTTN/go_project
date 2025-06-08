package handlers

import (
	"encoding/json"
	"game-matchmaking/backend/app"
	"game-matchmaking/backend/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func AddMapToGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		creatorID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var game models.Game
		if err := app.DB.Where("creator_id = ? AND status = ?", creatorID, "draft").First(&game).Error; err != nil {
			game = models.Game{CreatorID: creatorID.(uint), Status: "draft", CreatedAt: time.Now()}
			app.DB.Create(&game)
		}
		idStr := r.URL.Path[len("/api/games/add-map/"):]
		mapID, _ := strconv.Atoi(idStr)
		gameMap := models.GameMap{GameID: game.ID, MapID: uint(mapID)}
		app.DB.Create(&gameMap)
		app.DB.Preload("Maps").First(&game, game.ID)
		json.NewEncoder(w).Encode(game)
	}
}

func RemoveMapFromGame(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		if session.Values["user_id"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(r.URL.Path[len("/api/games/"):], "/remove-map/")
		gameID, _ := strconv.Atoi(parts[0])
		mapID, _ := strconv.Atoi(parts[1])
		app.DB.Where("game_id = ? AND map_id = ?", gameID, mapID).Delete(&models.GameMap{})
		json.NewEncoder(w).Encode(map[string]string{"message": "Map removed"})
	}
}
