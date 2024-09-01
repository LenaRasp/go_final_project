package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/utils"
)

func GetTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	utils.ResponseSuccess(w, task, http.StatusOK)
}
