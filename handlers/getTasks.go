package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/nextDate"
	"github.com/LenaRasp/go_final_project/response"
)

const limit = 10

func reqToDb(w http.ResponseWriter, rows *sql.Rows) []models.Task {
	var tasks []models.Task

	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			response.Error(w, "Ошибка сканирования БД", http.StatusInternalServerError)
			return nil
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		response.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return nil
	}

	if len(tasks) == 0 {
		tasks = make([]models.Task, 0, 0)
	}

	limitTasks := 10

	if len(tasks) > limitTasks {
		tasks = tasks[:limitTasks]
	}

	return tasks
}

func GetTasks(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var tasks []models.Task

	search := req.FormValue("search")

	if len(search) > 0 {
		dateTime, err := time.Parse("02.01.2006", search)
		if err != nil {
			rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit ", sql.Named("search", fmt.Sprint("%" + search + "%")), sql.Named("limit", limit))
			if err != nil {
				response.Error(w, "Ошибка запроса БД", http.StatusInternalServerError)
				return
			}

			tasks = reqToDb(w, rows)
		} else {
			rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit", sql.Named("date", dateTime.Format(nextDate.TimeLayout)), sql.Named("limit", limit))
			if err != nil {
				response.Error(w, "Ошибка запроса БД", http.StatusInternalServerError)
				return
			}

			tasks = reqToDb(w, rows)	
		}
	} else {
		rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
		if err != nil {
			response.Error(w, "Ошибка запроса БД", http.StatusInternalServerError)
			return
		}
		
		tasks = reqToDb(w, rows)
	}

	jsonResponse := map[string][]models.Task{"tasks": tasks}
	response.Success(w, jsonResponse, http.StatusOK)
}
