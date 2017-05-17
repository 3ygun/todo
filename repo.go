package main

import (
	"database/sql"
	"fmt"
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

	err := db.QueryRow("SELECT id, name, completed, due FROM "+dbTable+" WHERE id=?", id).Scan(&t.Id, &t.Name, &t.Completed, &t.Due)
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

	rows, err := db.Query("SELECT * FROM " + dbTable)
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
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (name, completed, due) VALUES(?, ?, ?)", dbTable))
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
	// for i, t := range todos {
	// 	if t.Id == id {
	// 		todos = append(todos[:i], todos[i+1:]...)
	// 		return nil
	// 	}
	// }
	// return fmt.Errorf("Could not find Todo with id of %d to delete", id)
	return nil
}
