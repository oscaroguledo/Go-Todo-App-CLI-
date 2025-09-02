package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

var tasks []Task
var taskFile = "tasks.json"

func addTask(description string) {
	newTask := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)
	fmt.Println("Task added:", description)
	saveTasks()

}
func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Println("Your Tasks:")
	for _, task := range tasks {
		status := "❌"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("%d. [%s] %s (Created at: %s)\n", task.ID, status, task.Description, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func deleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Println("Task deleted:", task.Description)
			saveTasks()
			return
		}
	}
	fmt.Println("Task not found with ID:", id)

}
func updateTask(id int, description string, completed bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].Completed = completed
			fmt.Println("Task updated:", description)
			saveTasks()
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}
func completeTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
			fmt.Println("Task marked as completed:", task.Description)
			saveTasks()
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}
func loadTasks() {
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return
		}
		fmt.Println("Error loading tasks:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("Error decoding tasks:", err)
	}
}

func saveTasks() {
	file, err := os.Create(taskFile)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding tasks:", err)
	}
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	loadTasks()
	fmt.Println("📝 Welcome to the CLI Todo App!")

	for {
		fmt.Println("\nChoose an option: add, list, update, complete, delete, exit")
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)

		switch input {
		case "add":

			fmt.Print("Enter task description: ")
			input, _ := reader.ReadString('\n') // Reads until Enter key
			descr := strings.TrimSpace(input)   // Removes newline at the end

			if descr == "" {
				fmt.Println("Task description cannot be empty.")
				continue
			}
			addTask(descr)
		case "list":

			listTasks()
		case "update":
			id := 0
			fmt.Print("Enter task ID to update: ")
			input, _ := reader.ReadString('\n') // Reads until Enter key
			input = strings.TrimSpace(input)    // Removes newline at the end

			id, err := strconv.Atoi(input) // convert string to int
			if err != nil {
				fmt.Println("Invalid ID. Please enter a number.")
				return
			}
			if id <= 0 {
				fmt.Println("Invalid task ID.")
				continue
			}
			fmt.Print("Enter new task description: ")
			descr, _ := reader.ReadString('\n') // Reads until Enter key
			descr = strings.TrimSpace(descr)    // Removes newline at the end

			fmt.Print("Is the task completed? (true/false): ")

			complete, _ := reader.ReadString('\n')
			complete = strings.TrimSpace(complete) // clean input
			completed, err := strconv.ParseBool(complete)
			if err != nil {
				fmt.Println("Invalid input. Please enter true or false.")
				return
			}
			updateTask(id, descr, completed)
		case "complete":
			fmt.Print("Enter task ID to complete: ")
			id := 0
			fmt.Print("Enter task ID to delete: ")
			fmt.Scanln(&id)
			if id <= 0 {
				fmt.Println("Invalid task ID.")
				continue
			}
			completeTask(id)
		case "delete":
			id := 0
			fmt.Print("Enter task ID to delete: ")
			fmt.Scanln(&id)
			if id <= 0 {
				fmt.Println("Invalid task ID.")
				continue
			}
			deleteTask(id)
		case "exit":
			saveTasks()
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
