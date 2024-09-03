package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"github.com/LenaRasp/go_final_project/database"
	"github.com/LenaRasp/go_final_project/handlers"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Ошибка при загрузке .env file")
	}
	PORT := os.Getenv("TODO_PORT")
	DBFILE := os.Getenv("TODO_DBFILE")

	db, err := sql.Open("sqlite", DBFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(appPath, DBFILE)
	_, err = os.Stat(dbFile)

	var install bool

	if err != nil {
		install = true
	}

	if install {
		database.CreateDB(db)
	} else {
		fmt.Println("База данных уже существует")
	}

	router := chi.NewRouter()

	router.Get("/*", handlers.GetStatic)
	router.Get("/api/nextdate", handlers.GetNextDate)
	router.Get("/api/tasks", func(w http.ResponseWriter, req *http.Request) { handlers.GetTasks(w, req, db) })
	router.Get("/api/task", func(w http.ResponseWriter, req *http.Request) { handlers.GetTask(w, req, db) })
	router.Put("/api/task", func(w http.ResponseWriter, req *http.Request) { handlers.UpdateTask(w, req, db) })
	router.Post("/api/task", func(w http.ResponseWriter, req *http.Request) { handlers.AddTask(w, req, db) })
	router.Post("/api/task/done", func(w http.ResponseWriter, req *http.Request) { handlers.DoneTask(w, req, db) })
	router.Delete("/api/task", func(w http.ResponseWriter, req *http.Request) { handlers.DeleteTask(w, req, db) })

	fmt.Println("Сервер прослушивает порт", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
