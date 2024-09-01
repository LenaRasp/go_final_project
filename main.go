package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LenaRasp/go_final_project/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Ошибка при загрузке .env file")
	}
	PORT := os.Getenv("TODO_PORT")

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool

	if err != nil {
		install = true
	}

	if install {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		)`)

		if err != nil {
			fmt.Println("Ошибка при создании таблицы:", err)
			return
		}

		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);")
		if err != nil {
			fmt.Println("Ошибка при создании индекса:", err)
			return
		}
		fmt.Println("Таблица и индекс успешно созданы")
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
