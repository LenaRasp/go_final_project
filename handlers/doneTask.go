package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/nextDate"
	"github.com/LenaRasp/go_final_project/utils"
)

func DoneTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	task := models.Task{}

	id := req.FormValue("id")
	if id == "" {
		utils.ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		utils.ResponseError(w, "Ошибка сканирования БД", http.StatusInternalServerError)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(utils.TimeLayout)
	}

	if task.Repeat != "" {
		newDate, err := nextDate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.ResponseError(w, "Ошибка формирования новой даты", http.StatusInternalServerError)
			return
		}
		task.Date = newDate

		_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("id", task.Id),
			sql.Named("date", task.Date))

		if err != nil {
			utils.ResponseError(w, "Ошибка db UPDATE", http.StatusInternalServerError)
			return
		}
	}
	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id))
		if err != nil {
			utils.ResponseError(w, "Ошибка db DELETE", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]models.Task{}
	utils.ResponseSuccess(w, response, http.StatusOK)
}
