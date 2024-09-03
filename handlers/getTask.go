package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/response"
)

func GetTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	task := models.Task{}

	id := req.FormValue("id")
	if id == "" {
		response.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		response.Error(w, "Ошибка сканирования БД", http.StatusInternalServerError)
		return
	}

	response.Success(w, task, http.StatusOK)
}
