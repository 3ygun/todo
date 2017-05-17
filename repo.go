package main

import (
	"database/sql"
	"log"
	"time"
)

// Give us some seed data
func CreateStartData() {
	RepoCreateTodo(Todo{
		Name:      "Write presentation",
		Completed: false,
		Due:       time.Now(),
	})
	RepoCreateTodo(Todo{
		Name:      "Host meetup",
		Completed: false,
		Due:       time.Now(),
	})
}

func RepoFindTodo(id int64) Todo {
	var t Todo

	err := db.QueryRow("SELECT id, name, completed, due "+
		"FROM "+dbTableData+", "+dbTableRemoved+" "+
		"WHERE id=? AND todos_id<>?", id, id).Scan(&t.Id, &t.Name, &t.Completed, &t.Due)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		return Todo{}
	case err != nil:
		log.Fatal(err)
	default:
		log.Printf("Todo is %s\n", t)
	}

	return t
}

func RepoGetAllTodos() []Todo {
	var todos []Todo

	rows, err := db.Query("SELECT * " +
		"FROM " + dbTableData + " " +
		"WHERE id NOT IN " +
		"(SELECT todo_id FROM " + dbTableRemoved + ");")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.Id, &t.Name, &t.Completed, &t.Due)
		if err != nil {
			panic(err)
		}
		todos = append(todos, t)
	}

	return todos
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(t Todo) Todo {
	stmt, err := db.Prepare("INSERT INTO " + dbTableData + " (name, completed, due) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(t.Name, t.Completed, t.Due)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	t.Id = lastId
	return t
}

func RepoDestroyTodo(id int64) error {
	stmt, err := db.Prepare("INSERT INTO " + dbTableRemoved + " (todo_id) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
