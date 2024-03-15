package main

import (
	"fmt"
	"net/http"

	"github.com/Crampustallin/todoList/templates"
	"github.com/a-h/templ"
)


func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mainPage := templates.Page()
	http.Handle("/", templ.Handler(mainPage))

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler(w,r)
			return
		}
		return
	})

	fmt.Println("Listening to :3000")
	http.ListenAndServe(":3000", nil)
}
