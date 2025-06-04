package main

import (
	"fmt"
	"log"
	"os"

	"tidyTodo/internal/common"
	"tidyTodo/internal/db"
)

func main() {
	DB, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	common.DB = DB

	if len(os.Args[1:]) > 0 {
		for _, title := range os.Args[1:] {
			_, err := DB.InsertTask(title)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var ch int

	for ch != -1 {
		tasks, err := DB.GetTasks()
		if err != nil {
			log.Fatal(err)
		}

		for index, task := range tasks {
			fmt.Println(index, task)
		}
		fmt.Print("> ")
		fmt.Scan(&ch)
		if ch < len(tasks) && ch > -1 {
			tsk := tasks[ch]
			err := DB.DeleteTask(tsk.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
