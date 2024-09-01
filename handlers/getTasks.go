package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/utils"
)

func GetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var tasks []models.Task

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC")
	if err != nil {
		utils.ResponseError(w, "Ошибка запроса БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			utils.ResponseError(w, "Ошибка сканирования БД", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		utils.ResponseError(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		tasks = make([]models.Task, 0, 0)
	}

	limitTasks := 10

	if len(tasks) > limitTasks {
		tasks = tasks[:limitTasks]
	}

	response := map[string][]models.Task{"tasks": tasks}
	utils.ResponseSuccess(w, response, http.StatusOK)
}
