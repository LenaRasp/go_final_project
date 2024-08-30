package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//Функция возвращает следующую дату в формате 20060102 и ошибку.

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Пустая строка")
	}
	//Возвращаемая дата должна быть больше даты, указанной в переменной now.
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	firstLetter := strings.Split(repeat, "")[0]
	if !(strings.EqualFold(firstLetter, "d") || strings.EqualFold(firstLetter, "y")) {
		return "", fmt.Errorf("Неподдерживаемый формат")
	}

	switch firstLetter {
	case "d":
		if len(repeat) < 2 {
			return "", fmt.Errorf("Неподдерживаемый формат")
		}
		newSlice := strings.Split(repeat, " ")

		days, err := strconv.Atoi(newSlice[1])

		if err != nil {
			fmt.Errorf("Неподдерживаемый формат")
		}

		if days < 1 || days > 400 {
			return "", fmt.Errorf("Неподдерживаемый формат")
		}

		dateTime = dateTime.AddDate(0, 0, days)
		for dateTime.Before(now) {
			dateTime = dateTime.AddDate(0, 0, days)
		}

	case "y":
		dateTime = dateTime.AddDate(1, 0, 0)
		for dateTime.Before(now) {
			dateTime = dateTime.AddDate(1, 0, 0)
		}
	}
	return dateTime.Format("20060102"), nil
}

func handleGetSite(w http.ResponseWriter, req *http.Request) {
	webDir := "web"

	if _, err := os.Stat(webDir + req.RequestURI); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, req)
	}
}

func handleGetNextDate(w http.ResponseWriter, req *http.Request) {
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	now, err := time.Parse("20060102", req.FormValue("now"))
	if err != nil {
		fmt.Errorf("Неподдерживаемый формат")
	}

	content, err := NextDate(now, date, repeat)
	if err != nil {
		fmt.Errorf("Неподдерживаемый формат")
	}

	fmt.Fprintf(w, content)
}

type Task struct {
	Id      string `json:"id"`      // id коллектива
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи. Обязательное поле
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения. Используется такой же формат, как в предыдущем шаге.
}

func ResponseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	errorMsg := map[string]string{"error": message}

	resp, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}

	w.Write(resp)
}

func ResponseSuccess(w http.ResponseWriter, response any, code int) {
	resp, err := json.Marshal(response)
	if err != nil {
		ResponseError(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func ValidateData(data Task) (Task, error) {
	task := data

	now := time.Now()

	if len(task.Title) == 0 {
		return task, fmt.Errorf("Не указано поле Заголовок")
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return task, fmt.Errorf("Неподдерживаемый формат даты")
	}

	if taskDate.Format("20060102") < now.Format("20060102") {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		}

		if task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return task, fmt.Errorf("Неподдерживаемый формат даты")
			}
			task.Date = nextDate
		}
	}

	return task, nil
}

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Ошибка при загрузке .env file")
	}
	PORT := os.Getenv("TODO_PORT")

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
	fmt.Println("install:", install)

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	router := chi.NewRouter()

	router.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		webDir := "web"

		if _, err := os.Stat(webDir + req.RequestURI); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.FileServer(http.Dir(webDir)).ServeHTTP(w, req)
		}
	})
	router.Get("/api/nextdate", handleGetNextDate)
	router.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		var tasks []Task

		rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC")
		if err != nil {
			ResponseError(w, "Ошибка запроса БД", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := Task{}

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				ResponseError(w, "Ошибка сканирования БД", http.StatusInternalServerError)
				return
			}

			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			log.Println(err)
			return
		}

		if len(tasks) == 0 {
			tasks = make([]Task, 0, 0)
		}

		limitTasks := 10

		if len(tasks) > limitTasks {
			tasks = tasks[:limitTasks]
		}

		response := map[string][]Task{"tasks": tasks}
		ResponseSuccess(w, response, http.StatusOK)
	})
	router.Get("/api/task", func(w http.ResponseWriter, req *http.Request) {
		id := req.FormValue("id")
		if id == "" {
			ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
			return
		}

		task := Task{}

		row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
			sql.Named("id", id))

		err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			ResponseError(w, "Ошибка сканирования БД", http.StatusInternalServerError)
			return
		}

		ResponseSuccess(w, task, http.StatusOK)
	})
	router.Put("/api/task", func(w http.ResponseWriter, req *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			ResponseError(w, "Ошибка ReadFrom", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			ResponseError(w, "Ошибка Unmarshal", http.StatusBadRequest)
			return
		}

		if task.Id == "" {
			ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
			return
		}

		_, err = strconv.Atoi(task.Id)
		if err != nil {
			ResponseError(w, "Невалидный идентификатор", http.StatusBadRequest)
			return
		}

		formattedTask, err := ValidateData(task)
		if err != nil {
			ResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		task = formattedTask

		res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("id", task.Id),
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))

		if err != nil {
			ResponseError(w, "Ошибка обновления db", http.StatusInternalServerError)
			return
		}

		row, err := res.RowsAffected()
		if err != nil {
			ResponseError(w, "Ошибка обновления db", http.StatusInternalServerError)
			return
		}
		if row == 0 {
			ResponseError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		}

		response := map[string][]Task{}
		ResponseSuccess(w, response, http.StatusOK)
	})
	router.Post("/api/task", func(w http.ResponseWriter, req *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			ResponseError(w, "Ошибка ReadFrom", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			ResponseError(w, "Ошибка Unmarshal", http.StatusBadRequest)
			return
		}

		formattedTask, err := ValidateData(task)
		if err != nil {
			ResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		task = formattedTask

		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))
		if err != nil {
			ResponseError(w, "Ошибка db", http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			ResponseError(w, "Ошибка id db", http.StatusInternalServerError)
			return
		}

		response := map[string]int64{"id": id}
		ResponseSuccess(w, response, http.StatusCreated)
	})
	router.Post("/api/task/done", func(w http.ResponseWriter, req *http.Request) {
		id := req.FormValue("id")
		if id == "" {
			ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
			return
		}

		task := Task{}

		row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
			sql.Named("id", id))

		err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			ResponseError(w, "Ошибка сканирования БД", http.StatusInternalServerError)
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat != "" {
			newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				ResponseError(w, "Ошибка формирования новой даты", http.StatusInternalServerError)
				return
			}
			task.Date = newDate

			_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
				sql.Named("id", task.Id),
				sql.Named("date", task.Date))

			if err != nil {
				ResponseError(w, "Ошибка db UPDATE", http.StatusInternalServerError)
				return
			}
		}
		if task.Repeat == "" {
			_, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
				sql.Named("id", id))
			if err != nil {
				ResponseError(w, "Ошибка db DELETE", http.StatusInternalServerError)
				return
			}
		}

		response := map[string]Task{}
		ResponseSuccess(w, response, http.StatusOK)
	})

	router.Delete("/api/task", func(w http.ResponseWriter, req *http.Request) {
		id := req.FormValue("id")
		if id == "" {
			ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
			return
		}

		res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id))
		if err != nil {
			ResponseError(w, "Ошибка db DELETE", http.StatusInternalServerError)
			return
		}
		
		row, err := res.RowsAffected()
		if err != nil {
			ResponseError(w, "Ошибка обновления db", http.StatusInternalServerError)
			return
		}
		if row == 0 {
			ResponseError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		}

		response := map[string]Task{}
		ResponseSuccess(w, response, http.StatusOK)
	})
	
	fmt.Println("Server is listening port", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
