package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Crampustallin/todoList/models"
	"github.com/Crampustallin/todoList/templates"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var todos map[int]models.Todo


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database", err)
	} else {
		fmt.Println("Connected to postgresql database")
	}

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	todos, err = getTodo(db)
	if err != nil {
		log.Fatal("Error getting todos:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			mainPage := templates.Page(todos)
			mainPage.Render(r.Context(), w)
		}
	})

	http.HandleFunc("/clicked/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			templates.Response(todos).Render(r.Context(), w)
			return
		}
		if r.Method == http.MethodPost {
			postHandler(w, r, db)
			return
		}
		if r.Method == http.MethodPut {
			putHandler(w,r, db)
			return
		}
		if r.Method == http.MethodDelete {
			delHandler(w,r, db)
			return 	
		}
		return
	})

	http.HandleFunc("/clicked/edit/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			editHandler(w,r)
			return
		}
	})

	fmt.Println("Listening to :3000")
	http.ListenAndServe(":3000", nil)
}
