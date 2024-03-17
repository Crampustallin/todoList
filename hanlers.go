package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Crampustallin/todoList/models"
	"github.com/Crampustallin/todoList/templates"
	_ "github.com/lib/pq"
)

func postHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	if r.Form.Has("todoDis") && r.Form.Has("status") {
		task := r.Form.Get("todoDis")
		status := r.Form.Get("status")
		saveTodo(db, task, status)
		res := templates.Response(todos)
		res.Render(r.Context(), w)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query().Get("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	todo := todos[id]
	templates.Edit(todo).Render(r.Context(), w)
}

func putHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := r.URL.Query().Get("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse from data", http.StatusBadRequest)
		return
	}
	if r.Form.Has("todoDis") && r.Form.Has("todoStatus") {
		task := r.Form.Get("todoDis")
		status := r.Form.Get("todoStatus")
		changeTodo(db, task, status, id)
	}
	templates.Response(todos).Render(r.Context(), w)
}

func delHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := r.URL.Query().Get("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = deleteTodo(db, id)
	if err != nil {
		http.Error(w, "Failed to delete from db", http.StatusInternalServerError)
		return
	}
	delete(todos, id)
	templates.Response(todos).Render(r.Context(), w)
}

func saveTodo(db *sql.DB, task, status string) error {
	id := 0
	err := db.QueryRow("INSERT INTO todo_list (description, status) VALUES ($1, $2) RETURNING task_id", task, status).Scan(&id)
	if err != nil {
		return fmt.Errorf("error inserting todo: %w", err)
	}
	todos[int(id)] = models.Todo{ID: int(id),  Description: task, Status: status}
	return nil
}

func changeTodo(db *sql.DB, task, status string, id int) error {
	_, err := db.Exec("UPDATE todo_list SET description=$1, status=$2 WHERE task_id=$3", task, status, id)
	if err != nil {
		return fmt.Errorf("error updating todolist: %w", err)
	}
	todos[id] = models.Todo{Description: task, Status: status}
	return nil
}

func getTodo(db *sql.DB) (map[int]models.Todo, error) {
	rows, err := db.Query("SELECT * FROM todo_list")
	if err != nil {
		return nil, fmt.Errorf("error querying todos: %w", err)
	}
	defer rows.Close()

	todos := make(map[int]models.Todo)
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Description, &todo.Status)
		if err != nil {
			return nil, fmt.Errorf("error scanning todo row: %w", err)
		}
		todos[todo.ID] = todo
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over todos rows: %w", err)
	}
	return todos, nil
}

func deleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM todo_list WHERE task_id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting from todolist: %w", err)
	}
	return nil
}
