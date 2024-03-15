package main

import (
	"net/http"
	"database/sql"
	"fmt"

	"github.com/Crampustallin/todoList/templates"
	"github.com/Crampustallin/todoList/models"
	_ "github.com/lib/pq"
)

func postHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}

	if r.Form.Has("todoDis") {
		task := r.Form.Get("todoDis")
		saveTodo(db, task, "in progress")
		todos = append(todos, models.Todo{Description: task})
		res := templates.Response(todos)
		res.Render(r.Context(), w)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}
	// TODO: crate update a todo in todo list
}

func saveTodo(db *sql.DB, task, status string) error {
	_, err := db.Exec("INSERT INTO todo_list (description, status) VALUES ($1, $2)", task, status)
	if err != nil {
		return fmt.Errorf("error inserting todo: %w", err)
	}
	return nil
}

func getTodo(db *sql.DB) ([]models.Todo, error) {
	rows, err := db.Query("SELECT * FROM todo_list")
	if err != nil {
		return nil, fmt.Errorf("error querying todos: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Description, &todo.Status)
		if err != nil {
			return nil, fmt.Errorf("error scanning todo row: %w", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over todos rows: %w", err)
	}
	return todos, nil
}
