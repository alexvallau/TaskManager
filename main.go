package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"news.com/events/models"
	"context"
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
	fmt.Println("The token is ", jwtToken)
	w.Header().Set("Authorization", jwtToken)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User logged in successfully")
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Zone test protégée OK")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			fmt.Printf("your token: %s", tokenString)
			http.Error(w, "Missing or invalid Authorization header {{tokenString}}", http.StatusUnauthorized)
		}
		//(fonctionne aussi)
		//tokenString = tokenString[len("Bearer "):]
		//tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		userid,err := models.VerifyToken(tokenString)
		//met la valeur de userid dans le contexte
		ctx := context.WithValue(r.Context(), "userId", userid)
		fmt.Printf("The user id from VerifyToken is %d", userid)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Invalid token. Could not access to the resource", http.StatusUnauthorized)
			return
		}
		if userid == -1{
			http.Error(w, "Missing information from User", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ShowProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(int)
	fmt.Println("Je suis dans ShowProjectHandler et l'id est ",userID)
	myTitleSlice := models.GetAllProject(userID)

	fmt.Fprintf(w,"%v",string(myTitleSlice))
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
	http.Handle("/createProject", AuthMiddleware(http.HandlerFunc(newProjectHandler)))
	http.Handle("/deleteProject", AuthMiddleware(http.HandlerFunc(deleteProjectHandler)))
	http.Handle("/project/createTask", AuthMiddleware(http.HandlerFunc(newTaskHandler)))
	http.Handle("/project/deleteTask", AuthMiddleware(http.HandlerFunc(deleteTaskHandler)))
	http.Handle("/test", AuthMiddleware(http.HandlerFunc(TestHandler)))
	http.Handle("/project/getAllProject", AuthMiddleware(http.HandlerFunc(ShowProjectHandler)))
	

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
