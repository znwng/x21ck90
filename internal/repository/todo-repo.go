package repository

import (
	"todo-go/internal/db"
	"todo-go/internal/models"
)

func GetAll() ([]models.Todo, error) {
	rows, err := db.DB.Query("SELECT id, item, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var todos []models.Todo

	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func GetByID(id int) (*models.Todo, error) {
	row := db.DB.QueryRow("SELECT id, item, completed FROM todos WHERE id = ?", id)

	var t models.Todo
	if err := row.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
		return nil, err
	}

	return &t, nil
}

func Create(item string) (int64, error) {
	res, err := db.DB.Exec(
		"INSERT INTO todos (item, completed) VALUES (?, ?)",
		item,
		false,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Toggle(id int) error {
	_, err := db.DB.Exec(
		"UPDATE todos SET completed = NOT completed WHERE id = ?",
		id,
	)
	return err
}

