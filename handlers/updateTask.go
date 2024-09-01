package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/utils"
	"github.com/LenaRasp/go_final_project/validations"
)

func UpdateTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		utils.ResponseError(w, "Ошибка ReadFrom", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		utils.ResponseError(w, "Ошибка Unmarshal", http.StatusBadRequest)
		return
	}

	if task.Id == "" {
		utils.ResponseError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(task.Id)
	if err != nil {
		utils.ResponseError(w, "Невалидный идентификатор", http.StatusBadRequest)
		return
	}

	formattedTask, err := validations.ValidateData(task)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
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
		utils.ResponseError(w, "Ошибка обновления db", http.StatusInternalServerError)
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

	response := map[string][]models.Task{}
	utils.ResponseSuccess(w, response, http.StatusOK)
}
