package main

import (
	"fmt"
	"net/http"
	"strconv"

	"shopls/logging"
	"shopls/templates"
	"shopls/types"

	"github.com/go-chi/chi/v5"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// CONFIG
var DEFAULT_TITLE = "shop-ls"
var PORT = ":8080"

// GLOBAL STATE
var DATABASE *sql.DB

func handle404(w http.ResponseWriter, r *http.Request) {
	component := templates.Layout(templates.Title(DEFAULT_TITLE), templates.Error404())
	component.Render(r.Context(), w)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handle404(w, r)
	} else {
		component := templates.Layout(templates.Title(DEFAULT_TITLE), templates.Index())
		component.Render(r.Context(), w)
	}
}

func handleItemsGet(w http.ResponseWriter, r *http.Request) {
	rows, err := DATABASE.Query("SELECT id, description FROM items")
	if err != nil {
		logging.Error("Error querying database: %s", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []types.Item{}
	for rows.Next() {
		var item types.Item
		err := rows.Scan(&item.Id, &item.Description)
		if err != nil {
			logging.Error("Error scanning rows: %s", err)
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	component := templates.Items(items)
	component.Render(r.Context(), w)
}

func handleItemsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	description := r.FormValue("description")
	res, err := DATABASE.Exec("INSERT INTO items (description) VALUES (?)", description)
	if err != nil {
		logging.Error("Error inserting into database: %s", err)
		http.Error(w, "Error inserting into database", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		logging.Error("Error getting last insert id: %s", err)
		http.Error(w, "Error getting last insert id", http.StatusInternalServerError)
		return
	}

	component := templates.Item(types.Item{Description: description, Id: id})
	component.Render(r.Context(), w)
}

func handleItemsPut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		logging.Error("Error parsing id: %s", err)
		http.Error(w, "Invalid request, id not found", http.StatusBadRequest)
		return
	}

	_, err = DATABASE.Exec("UPDATE items SET description = ? WHERE id = ?", r.FormValue("description"), id)
	if err != nil {
		logging.Error("Error updating database: %s", err)
		http.Error(w, "Error updating database", http.StatusInternalServerError)
		return
	}

	component := templates.Item(types.Item{Description: r.FormValue("description"), Id: id})
	component.Render(r.Context(), w)
}

func handleItemsDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		logging.Error("Error parsing id: %s", err)
		http.Error(w, "Invalid request, id not found", http.StatusBadRequest)
		return
	}

	_, err = DATABASE.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		logging.Error("Error deleting from database: %s", err)
		http.Error(w, "Error deleting from database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "")
}

func handleItemsEditor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		logging.Error("Error parsing id: %s", err)
		http.Error(w, "Invalid request, id not found", http.StatusBadRequest)
		return
	}
	component := templates.ItemEdit(types.Item{Description: r.FormValue("description"), Id: id})
	component.Render(r.Context(), w)
}

func main() {
	// LOGGING
	logging.Init("", true)

	// DATABASE
	var err error
	DATABASE, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		logging.Fatal("Error opening database: %s", err)
	}
	defer DATABASE.Close()

	// TABLE
	_, err = DATABASE.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, description TEXT)")

	// HANDLERS
	r := chi.NewRouter()
	r.Get("/", handleIndex)
	r.Get("/items", handleItemsGet)
	r.Post("/items", handleItemsPost)
	r.Put("/items/{id}", handleItemsPut)
	r.Delete("/items/{id}", handleItemsDelete)
	r.Post("/items/editor", handleItemsEditor)

	// SERVER
	logging.Info("Server started at port %s", PORT)
	http.ListenAndServe(PORT, r)
}
