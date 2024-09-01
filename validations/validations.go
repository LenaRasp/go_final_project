package validations

import (
	"fmt"
	"time"

	"github.com/LenaRasp/go_final_project/models"
	"github.com/LenaRasp/go_final_project/nextDate"
	"github.com/LenaRasp/go_final_project/utils"
)

func ValidateData(data models.Task) (models.Task, error) {
	task := data

	now := time.Now()

	if len(task.Title) == 0 {
		return task, fmt.Errorf("Не указано поле Заголовок")
	}

	if task.Date == "" {
		task.Date = now.Format(utils.TimeLayout)
	}

	taskDate, err := time.Parse(utils.TimeLayout, task.Date)
	if err != nil {
		return task, fmt.Errorf("Неподдерживаемый формат даты")
	}

	if taskDate.Format(utils.TimeLayout) < now.Format(utils.TimeLayout) {
		if task.Repeat == "" {
			task.Date = now.Format(utils.TimeLayout)
		}

		if task.Repeat != "" {
			nextDate, err := nextDate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return task, fmt.Errorf("Неподдерживаемый формат даты")
			}
			task.Date = nextDate
		}
	}

	return task, nil
}
