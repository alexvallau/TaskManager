package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type State string

const (
	InProgress State = "in progress"
	Done       State = "done"
)

type Task struct {
	Id          int    `json:"id,omitempty"`
	ProjectId   int    `json:"projectid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       State  `json:"state"`
	Comment     string `json:"comment"`
}

type Project struct {
	Id          int    `json:"id,omitempty"`
	OwnerId     int    `json:"ownerId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       State  `json:"state"`
}

type Utilisateur struct {
	JWT      string `json:"jwt,omitempty"`
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

func (p *Project) CreateProject() {

	db, err := connectDB()
	if err != nil {
		log.Panic(err)
		fmt.Println("Could not connect to BDD")
	}
	defer db.Close()
	query := "INSERT INTO projects(owner_id, title, Description) VALUES (?,?, ?)"
	_, err = db.Exec(query, p.OwnerId, p.Title, p.Description)
	if err != nil {
		log.Panic(err)
		fmt.Println("Could not insert Project into BDD")
	}
	fmt.Printf(" %s  was added", p.Title)
}

func (p *Project) DeleteProject() {

	db, err := connectDB()
	if err != nil {
		log.Panic(err)
		fmt.Println("Could not connect to BDD")
	}
	defer db.Close()
	//check si existe
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM projetcs WHERE id = ?)"
	err = db.QueryRow(checkQuery, p.Id).Scan(&exists)
	if err != nil{
		log.Panic(err)
		fmt.Println("Error Checking if project exists")
	}
	if !exists{
		fmt.Printf("Le projet n'existe pas")
		return
	}
	//Si le projet n'existe pas, on supprime
	query := "DELETE FROM projects WHERE id = ?"
	_, err = db.Exec(query, p.Id)
	if err != nil {
		log.Panic(err)
		fmt.Println("Could not insert Project into BDD")
	}
	fmt.Printf(" Project number %d  was deleted", p.Id)
}

func (t *Task) CreateTask() {
	db, err := connectDB()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	query := "INSERT INTO tasks(project_id, title, Description, state, commentaire) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, t.ProjectId, t.Title, t.Description, t.State, t.Comment)
	if err != nil {
		log.Panic(err)
		fmt.Println("Could not insert task into BDD")
	}
	fmt.Printf("Task %s  was added", t.Title)
}

func (t *Task) DeleteTask() {
	db, err := connectDB()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	query := "DELETE FROM tasks WHERE id = ?"
	_, err = db.Exec(query, t.Id)
	if err != nil {
		log.Panic(err)
		fmt.Printf("Could not delete task %d", t.Id)
	}
	fmt.Printf("Task %d  was deleted", t.Id)
}
