package main

import (
	"net/http"
	"github.com/Crampustallin/todoList/templates"
)


func postHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}

	if r.Form.Has("todoDis") {
		res := templates.Response(r.Form.Get("todoDis"))
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
