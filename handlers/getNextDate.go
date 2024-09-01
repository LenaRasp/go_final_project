package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LenaRasp/go_final_project/nextDate"
	"github.com/LenaRasp/go_final_project/utils"
)

func GetNextDate(w http.ResponseWriter, req *http.Request) {
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	now, err := time.Parse(utils.TimeLayout, req.FormValue("now"))
	if err != nil {
		fmt.Errorf("Неподдерживаемый формат")
	}

	content, err := nextDate.NextDate(now, date, repeat)
	if err != nil {
		fmt.Errorf("Неподдерживаемый формат")
	}

	fmt.Fprintf(w, content)
}
