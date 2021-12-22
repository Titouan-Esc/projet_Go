package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ! Création struc d'une personne
type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

// ! Création struc d'un livre
type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var (
	person = &Person{
		Name:  "Titouan",
		Email: "titouan@outlook.fr",
	}
	books = []Book{
		{
			Title:      "La vie trépidente de Titouan",
			Author:     "Titouan Escorneboueu",
			CallNumber: 1234,
			PersonID:   1,
		},
		{
			Title:      "Les histoires de l'univers",
			Author:     "Titouan Escorneboueu",
			CallNumber: 5678,
			PersonID:   1,
		},
	}
)

var db *gorm.DB
var err error

func main() {
	// * Récupérer les variables d'environnement
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// * Connection à la base de donnée
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	// * Ouvrir la connection à la base de donnée
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("La connection à la BDD est correcte")
	}

	// * Fermer la base de donnée lorsque la fonction main est fini
	postgresDB, err := db.DB()
	if err != nil {
		log.Fatal()
	}
	defer postgresDB.Close()

	// * Créer la migration vers la base de donnée si elle n'a pas déjà était créé
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// todo (simple formule pour créer, m'a servis à savoir si tout marcher bien)
	// db.Create(person)
	// for idx := range books {
	// 	db.Create(&books[idx])
	// }

	// * Routes de l'API
	router := mux.NewRouter()

	// ? GET
	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET") // Ainsi que les livres
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")

	// ? POST
	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/create/book", createBook).Methods("POST")

	// ? DELETE
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	router.HandleFunc("/delete/book/{id}", deleteBook).Methods("DELETE")

	// ? PORT
	http.ListenAndServe(":8000", router)
}

// ! API controllers

// TODO | Person controllers
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)
	json.NewEncoder(w).Encode(&people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	var books []Book

	db.First(&person, params["id"])
	db.Model(&person).Association("Books").Find(&books) // Méthodes pour récupérer les books qui sont dans le mode Person

	person.Books = books

	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}

// TODO | Books controllers
func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(&book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}
