package handlers

import (
	"encoding/json"
	"game-matchmaking/backend/app"
	"game-matchmaking/backend/models"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var pool = &redis.Pool{
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

func RoleBasedMiddleware(next http.HandlerFunc, allowedRoles []string, app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := app.Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		userID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var player models.Player
		if err := app.DB.First(&player, userID).Error; err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		hasPermission := false
		for _, role := range allowedRoles {
			if player.Role == role {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func Register(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if input.Username == "" || input.Email == "" || input.Password == "" || input.Role == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		player := models.Player{
			Username: input.Username,
			Email:    input.Email,
			Password: string(hashedPassword),
			Role:     input.Role,
			Token:    uuid.New().String(),
		}
		if err := app.DB.Create(&player).Error; err != nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		session, _ := app.Store.Get(r, "session")
		session.Values["user_id"] = player.ID
		session.Values["token"] = player.Token
		session.Save(r, w)
		json.NewEncoder(w).Encode(player)
	}
}

func Login(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if input.Username == "" || input.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}
		var dbPlayer models.Player
		if err := app.DB.Where("username = ?", input.Username).First(&dbPlayer).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(dbPlayer.Password), []byte(input.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if dbPlayer.Token == "" {
			dbPlayer.Token = uuid.New().String()
			app.DB.Save(&dbPlayer)
		}
		session, err := app.Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = dbPlayer.ID
		session.Values["token"] = dbPlayer.Token
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}
		conn := pool.Get()
		defer conn.Close()
		conn.Do("SET", "session:"+dbPlayer.Username, dbPlayer.ID, "EX", 86400*7)
		json.NewEncoder(w).Encode(dbPlayer)
	}
}

func ValidateSession(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := app.Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		userID, ok := session.Values["user_id"]
		token, tokenOk := session.Values["token"]
		if !ok || !tokenOk {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var player models.Player
		if err := app.DB.First(&player, userID).Error; err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if player.Token != token {
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}
		conn := pool.Get()
		defer conn.Close()
		_, err = conn.Do("GET", "session:"+player.Username)
		if err != nil {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(player)
	}
}

func Logout(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if ok {
			var player models.Player
			if err := app.DB.First(&player, userID).Error; err == nil {
				player.Token = ""
				app.DB.Save(&player)
			}
			conn := pool.Get()
			defer conn.Close()
			conn.Do("DEL", "session:"+player.Username)
		}
		session.Values["user_id"] = nil
		session.Values["token"] = nil
		session.Save(r, w)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
	}
}

func ManageUsers(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var players []models.Player
		app.DB.Find(&players)
		json.NewEncoder(w).Encode(players)
	}
}
