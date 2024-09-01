package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/utils"
)

func DeleteTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id := req.FormValue("id")
	if id == "" {
		utils.ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		utils.ResponseError(w, "Ошибка db DELETE", http.StatusInternalServerError)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		utils.ResponseError(w, "Ошибка обновления db", http.StatusInternalServerError)
		return
	}
	if row == 0 {
		utils.ResponseError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	response := map[string]models.Task{}
	utils.ResponseSuccess(w, response, http.StatusOK)
}
