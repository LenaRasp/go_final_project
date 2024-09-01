package models

type Task struct {
	Id      string `json:"id"`      // id коллектива
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи. Обязательное поле
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения. Используется такой же формат, как в предыдущем шаге.
}
