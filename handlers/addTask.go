package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/utils"
	"github.com/LenaRasp/go_final_project/validations"
)

func AddTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	formattedTask, err := validations.ValidateData(task)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task = formattedTask

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		utils.ResponseError(w, "Ошибка db", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		utils.ResponseError(w, "Ошибка id db", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"id": id}
	utils.ResponseSuccess(w, response, http.StatusCreated)
}
