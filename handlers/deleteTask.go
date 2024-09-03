package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/response"
)

func DeleteTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id := req.FormValue("id")
	if id == "" {
		response.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		response.Error(w, "Ошибка db DELETE", http.StatusInternalServerError)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		response.Error(w, "Ошибка обновления db", http.StatusInternalServerError)
		return
	}
	if row == 0 {
		response.Error(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]models.Task{}
	response.Success(w, jsonResponse, http.StatusOK)
}
