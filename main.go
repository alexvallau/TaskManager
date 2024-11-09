package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	isLoggedIn, err := user.Login()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isLoggedIn {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//jwtToken
	if user.Id == -1 {
		return
	}
	jwtToken, err := user.GenerateJWT(user.Id)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	//response := map[string]string{
	//	"jwt": jwtToken,
	//}
	w.Header().Set("Authorization", "Bearer "+jwtToken)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User logged in successfully")
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := models.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bravo")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		fmt.Printf("hi x2 %s",tokenString)
		if tokenString == "" {
			fmt.Printf("your token: %s", tokenString)
			http.Error(w, "Missing or invalid Authorization header {{tokenString}}", http.StatusUnauthorized)
		}

		//(fonctionne aussi)
		//tokenString = tokenString[len("Bearer "):] 
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		err := models.VerifyToken(tokenString)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Invalid token. Could not access to the resource", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
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
	http.Handle("/createProject",  AuthMiddleware(http.HandlerFunc(newProjectHandler)))
	http.Handle("/deleteProject",  AuthMiddleware(http.HandlerFunc(deleteProjectHandler)))
	http.Handle("/createTask",  AuthMiddleware(http.HandlerFunc(newTaskHandler)))
	http.Handle("/deleteTask",  AuthMiddleware(http.HandlerFunc(deleteTaskHandler)))
	http.Handle("/test", AuthMiddleware(http.HandlerFunc(TestHandler)))

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
