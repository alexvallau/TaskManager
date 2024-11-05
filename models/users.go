package models

import (
	
	"database/sql"
	"fmt"
	"log"
	"encoding/json"
	"fmt"
	
	"log"
	
	"net/http"
)
	"golang.org/x/crypto/bcrypt"
)

type Utilisateur struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func connectDB() (*sql.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/tasksapi"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func hashPassword(UnhashedPassword string) (hashedPassword []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(UnhashedPassword), bcrypt.DefaultCost)
}

func CreateUser(username, password string) {

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hashedPassword, err := hashPassword(password)

	if err != nil {
		log.Fatal(err)
	}

	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = db.Exec(query, username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User admin was added ")
}

func (u *Utilisateur) Login() (bool, error) {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()
	var hashedPassword string
	query := "SELECT password FROM users WHERE username = ?"
	err = db.QueryRow(query, u.Username).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	if err != nil {
		fmt.Printf("Failed Connexion with username %s and password %s \n", u.Username, u.Password)
		return false, err
	}
	fmt.Printf("User %s correctly logged in", u.Username)
	return true, nil
}

type Tache struct {
	Titre        string
	Description  string
	Priorité     string
	Etat         string
	Utilisateurs []Utilisateur
}

type Projet struct {
	Titre        string
	Etat         string
	Taches       []Tache
	Utilisateurs []Utilisateur
}


