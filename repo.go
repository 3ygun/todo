package main

import (
	"fmt"
	"log"
	"time"
)

var currentId int

var todos Todos

// Give us some seed data
func CreateStartData() {
	RepoCreateTodo("Write presentation")
	RepoCreateTodo("Host meetup")
}

func RepoFindTodo(id int64) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(name string) Todo {
	datetime := time.Now()
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (name, completed, due) VALUES(?, FALSE, NOW())", dbTable))
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	t := Todo{
		Id:        lastId,
		Name:      name,
		Completed: false,
		Due:       datetime,
	}
	todos = append(todos, t)
	return t
}

func RepoCreateTodoFrom(todo Todo) Todo {
	todo.Id = 0
	todos = append(todos, todo)
	return todo
}

func RepoDestroyTodo(id int64) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
