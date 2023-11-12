// main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Customer struct {
	gorm.Model
	Username   string
	Email      string
	Password   string
	CusAddress string
	CusPh      string
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/sign_up", signUpHandler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/sign.html", http.StatusSeeOther)
	}).Methods("GET")

	port := 3020
	fmt.Printf("Listening on PORT %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func initDB() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db") // Change to your preferred database connection details
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Customer{})
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	address := r.FormValue("address")
	phno := r.FormValue("phno")

	customer := Customer{
		Username:   name,
		Email:      email,
		Password:   password,
		CusAddress: address,
		CusPh:      phno,
	}

	if err := db.Create(&customer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Record Inserted Successfully")
	http.Redirect(w, r, "/home.html", http.StatusSeeOther)
}
