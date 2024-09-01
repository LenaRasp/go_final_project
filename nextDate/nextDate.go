package nextDate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LenaRasp/go_final_project/utils"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Пустая строка")
	}
	//Возвращаемая дата должна быть больше даты, указанной в переменной now.
	dateTime, err := time.Parse(utils.TimeLayout, date)
	if err != nil {
		return "", err
	}

	firstLetter := strings.Split(repeat, "")[0]
	if !(strings.EqualFold(firstLetter, "d") || strings.EqualFold(firstLetter, "y")) {
		return "", fmt.Errorf("Неподдерживаемый формат")
	}

	switch firstLetter {
	case "d":
		if len(repeat) < 2 {
			return "", fmt.Errorf("Неподдерживаемый формат")
		}
		newSlice := strings.Split(repeat, " ")

		days, err := strconv.Atoi(newSlice[1])

		if err != nil {
			fmt.Errorf("Неподдерживаемый формат")
		}

		if days < 1 || days > 400 {
			return "", fmt.Errorf("Неподдерживаемый формат")
		}

		dateTime = dateTime.AddDate(0, 0, days)
		for dateTime.Before(now) {
			dateTime = dateTime.AddDate(0, 0, days)
		}

	case "y":
		dateTime = dateTime.AddDate(1, 0, 0)
		for dateTime.Before(now) {
			dateTime = dateTime.AddDate(1, 0, 0)
		}
	}
	return dateTime.Format(utils.TimeLayout), nil
}
