package main
// underscore is a special character in Go that tells the compiler to import the package but not use it in the code.
import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Task represents a to-do item
// struct is a keyword in Go that declares a data structure it's like a class in other languages
type Task struct {
	ID   int
	Name string
}

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create tasks table if it doesn't exist
	//the underscore (_) is used to discard the first value returned by the db.Exec() function call. In this case the first value returned is the result of the query, which we don't need.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		fmt.Println("Error creating table:", err)
		os.Exit(1)
	}

	// Print welcome message
	fmt.Println("Welcome to the to-do list manager!")
	fmt.Println("Type 'add' to add a task, 'delete' to delete a task, or 'view' to view all tasks.")

	// Loop to read user input and execute commands
	for {
		fmt.Print("> ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "add":
			fmt.Println("Enter the name of the task:")
			var name string
			fmt.Scanln(&name)
			err = addTask(db, name)
			if err != nil {
				fmt.Println("Error adding task:", err)
			} else {
				fmt.Println("Task added successfully!")
			}
		case "delete":
			fmt.Println("Enter the ID of the task to delete:")
			var id int
			fmt.Scanln(&id)
			err = deleteTask(db, id)
			if err != nil {
				fmt.Println("Error deleting task:", err)
			} else {
				fmt.Println("Task deleted successfully!")
			}
		case "view":
			tasks, err := getTasks(db)
			if err != nil {
				fmt.Println("Error getting tasks:", err)
			} else {
				fmt.Println("Tasks:")
				for _, task := range tasks {
					fmt.Printf("%d. %s\n", task.ID, task.Name)
				}
			}
		default:
			fmt.Println("Invalid command. Type 'add', 'delete', or 'view'.")
		}
	}
}

// addTask adds a task to the database
// the * is a pointer in Go that points to the memory address of the variable
func addTask(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO tasks (name) VALUES (?)", name)
	return err
}

// deleteTask deletes a task from the database by ID
func deleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// getTasks gets all tasks from the database
func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
