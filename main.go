package main

import (
	"fmt"
	"net/http"

	"github.com/Crampustallin/todoList/templates"
)

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	fmt.Println("Listening to :3000")
	http.ListenAndServe(":3000", nil)
}
