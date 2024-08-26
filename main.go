package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"database/sql"
	// "bytes"
	"encoding/json"
	_ "modernc.org/sqlite"


	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
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
		fmt.Println(date, ">", dateTime.Format("20060102"), repeat, days)
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
	id            string        `json:"id"`     // id коллектива
	date          string        `json:"date"`             // дата задачи в формате 20060102
	title				 	string      	`json:"title"`    // заголовок задачи. Обязательное поле
	comment       string        `json:"comment"`           // комментарий к задаче
	repeat			  string        `json:"repeat"`    // правило повторения. Используется такой же формат, как в предыдущем шаге.
}

// func postTasks(w http.ResponseWriter, req *http.Request) {
// 	var task Task
// 	var buf bytes.Buffer

// 	_, err := buf.ReadFrom(req.Body)
// 	if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 	}

// 	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 	}
// 	// showCode := &task
// 	fmt.Println(task, task.title, task.date)
// 	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)", 
// 		sql.Named("date", task.date),
// 		sql.Named("title", task.title),
// 		sql.Named("comment", task.comment),
// 		sql.Named("repeat", task.repeat))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	res.LastInsertId()

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// }

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
	router.Post("/api/task", func(w http.ResponseWriter, req *http.Request) {
		var task Task
		// var buf bytes.Buffer
	
		// _, err := buf.ReadFrom(req.Body)
		// if err != nil {
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// }
	
		// if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// }

		err := json.NewDecoder(req.Body).Decode(&task)
		if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
		}

		taskDate, err := time.Parse("20060102", task.date)
		fmt.Println(task.date, taskDate, err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.date == "" {
			task.date = time.Now().Format("20060102")
		}

		if len(task.title) == 0 {
			http.Error(w, "Не указано поле title", http.StatusBadRequest)
			return
		}

		if taskDate.Before(time.Now()) {
			if task.repeat == "" {
				task.date = time.Now().Format("20060102")
			}
			//при указанном правиле повторения вам нужно вычислить и записать в таблицу дату выполнения, 
			//которая будет больше сегодняшнего числа. 
			//Для этого используйте функцию NextDate(), которую вы уже написали раньше.
			if task.repeat != "" {
				task.date, err = NextDate(time.Now(), task.date, task.repeat)
			} 
		}
		

		fmt.Println(task, task.title, task.date)
		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)", 
			sql.Named("date", task.date),
			sql.Named("title", task.title),
			sql.Named("comment", task.comment),
			sql.Named("repeat", task.repeat))
		if err != nil {
			fmt.Println(err)
			return
		}
	
	id, err := res.LastInsertId()
	fmt.Println(res.LastInsertId())
    
	mes := map[string]int64{"id": id}
	resp, err := json.Marshal(mes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	})
	// запускаем сервер
	fmt.Println("Server is listening port", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
