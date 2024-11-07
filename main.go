package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"news.com/events/models"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	//Vérifie que la méthode est bien un poste
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	//Décode depuis le format json et instancie l'utilisateur
	var user models.Utilisateur
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
		return
	}
	var userId int
	isLoggedIn, err := user.Login()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isLoggedIn {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "User %s with id %d logged in successfully \n", user.Username, userId)

	//jwtToken
	if user.Id == -1{
		return
	}
	jwtToken, err := user.GenerateJWT(user.Id)
	if err != nil {
		http.Error(w,"Could not generate token", http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"jwt": jwtToken,
	}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method in newProjectHandler", http.StatusMethodNotAllowed)
	}
	var projet models.Project
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&projet)
	if err != nil {
		log.Fatal(err)
		return
	}
	projet.CreateProject()
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method in newProjectHandler", http.StatusMethodNotAllowed)
	}
	var projet models.Project
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&projet)
	if err != nil {
		log.Fatal(err)
		return
	}
	projet.DeleteProject()
}

func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method in newTasktHandler", http.StatusMethodNotAllowed)
	}
	var task models.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		log.Fatal(err)
		return
	}
	task.CreateTask()
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method in newTasktHandler", http.StatusMethodNotAllowed)
	}
	var task models.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		log.Fatal(err)
		return
	}
	task.DeleteTask()
}

func main() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/createProject", newProjectHandler)
	http.HandleFunc("/deleteProject", deleteProjectHandler)
	http.HandleFunc("/createTask", newTaskHandler)
	http.HandleFunc("/deleteTask", deleteTaskHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
