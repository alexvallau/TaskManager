package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"news.com/events/models"

	_ "github.com/go-sql-driver/mysql"
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
	fmt.Fprintf(w, "User %s logged in successfully \n", user.Username)
}

func main() {

	models.CreateUser("user", "user")
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
