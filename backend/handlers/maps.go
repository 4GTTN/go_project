package handlers

import (
	"backend/app"
	"backend/models"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func ListHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var maps []models.Map
		filter := r.URL.Query().Get("filter")
		query := app.DB.Where("is_deleted = ?", false)
		if filter != "" {
			query = query.Where("name ILIKE ? OR CAST(price AS TEXT) ILIKE ?", "%"+filter+"%", "%"+filter+"%")
		}
		query.Find(&maps)
		data := struct {
			Maps   []models.Map
			Filter string
		}{maps, filter}
		templates.ExecuteTemplate(w, "list.html", data)
	}
}

func DetailHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/map/"):]
		id, _ := strconv.Atoi(idStr)
		var mapObj models.Map
		if err := app.DB.Where("id = ? AND is_deleted = ?", id, false).First(&mapObj).Error; err != nil {
			http.NotFound(w, r)
			return
		}
		templates.ExecuteTemplate(w, "detail.html", struct{ Map models.Map }{mapObj})
	}
}

func DeleteHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/delete/"):]
		id, _ := strconv.Atoi(idStr)
		app.DB.Exec("UPDATE maps SET is_deleted = true WHERE id = ?", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetMaps(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var maps []models.Map
		filter := r.URL.Query().Get("filter")
		query := app.DB.Where("is_deleted = ?", false)
		if filter != "" {
			query = query.Where("name ILIKE ? OR CAST(price AS TEXT) ILIKE ?", "%"+filter+"%", "%"+filter+"%")
		}
		query.Find(&maps)
		json.NewEncoder(w).Encode(maps)
	}
}

func GetMap(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/maps/"):]
		id, _ := strconv.Atoi(idStr)
		var mapObj models.Map
		if err := app.DB.Where("id = ? AND is_deleted = ?", id, false).First(&mapObj).Error; err != nil {
			http.Error(w, "Map not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(mapObj)
	}
}

func CreateMap(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.Store.Get(r, "session")
		if session.Values["user_id"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var mapObj models.Map
		if err := json.NewDecoder(r.Body).Decode(&mapObj); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		app.DB.Create(&mapObj)
		json.NewEncoder(w).Encode(mapObj)
	}
}
