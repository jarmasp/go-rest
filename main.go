package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTask []Task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Go Lang",
		Content: "Create rest API using Go",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func newTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "insert a valid task")
		return
	}
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID is not a number")
		return
	}
	for _, task := range tasks {
		if task.ID == taskId {
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID is not a number")
		return
	}
	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Fprintf(w, "deleted task with id %v", taskId)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	var updatedTask Task
	if err != nil {
		fmt.Fprintf(w, "ID is not a number")
		return
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "insert a valid task")
		return
	}
	json.Unmarshal(reqBody, &updatedTask)
	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			updatedTask.ID = taskId
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "updated task with id %v", taskId)
		}
	}
}

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", IndexRoute).Methods("GET")
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/task/{id}", getTask).Methods("GET")
	router.HandleFunc("/task", newTask).Methods("POST")
	router.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", router))
}
