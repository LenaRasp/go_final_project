package database

import (
	"database/sql"
	"fmt"
)

func CreateDB(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		)`)

	if err != nil {
		fmt.Println("Ошибка при создании таблицы:", err)
		return
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);")
	if err != nil {
		fmt.Println("Ошибка при создании индекса:", err)
		return
	}
	fmt.Println("Таблица и индекс успешно созданы")
}
