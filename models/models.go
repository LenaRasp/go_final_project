package models

import (
	"fmt"
	"time"

	"github.com/LenaRasp/go_final_project/nextDate"

)

type Task struct {
	Id      string `json:"id"`      // id коллектива
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи. Обязательное поле
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения. Используется такой же формат, как в предыдущем шаге.
}

func (t Task) ValidateData() (Task, error) {
	now := time.Now()

	if len(t.Title) == 0 {
		return t, fmt.Errorf("Не указано поле Заголовок")
	}

	if t.Date == "" {
		t.Date = now.Format(nextDate.TimeLayout)
	}

	taskDate, err := time.Parse(nextDate.TimeLayout, t.Date)
	if err != nil {
		return t, fmt.Errorf("Неподдерживаемый формат даты")
	}

	if taskDate.Format(nextDate.TimeLayout) < now.Format(nextDate.TimeLayout) {
		if t.Repeat == "" {
			t.Date = now.Format(nextDate.TimeLayout)
		}

		if t.Repeat != "" {
			nextDate, err := nextDate.NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return t, fmt.Errorf("Неподдерживаемый формат даты")
			}
			t.Date = nextDate
		}
	}

	return t, nil
}
