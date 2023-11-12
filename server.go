// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

var session *gocql.Session

// Customer struct represents the schema for customers in Cassandra
type Customer struct {
	Username   string
	Email      string
	Password   string
	CusAddress string
	CusPh      string
}

func init() {
	// Connect to Cassandra
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "coffee"
	session, _ = cluster.CreateSession()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/loginDetails", loginDetailsHandler).Methods("POST")
	// Add other routes...

	port := ":9998"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "log.html", http.StatusSeeOther)
}

func loginDetailsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("userName")
	password := r.FormValue("passw")
	checkPass(w, r, email, password)
}

func checkPass(w http.ResponseWriter, r *http.Request, email, password string) {
	var customer Customer
	query := "SELECT * FROM customers WHERE email=? AND password=? LIMIT 1"

	if err := session.Query(query, email, password).Consistency(gocql.One).Scan(
		&customer.Username,
		&customer.Email,
		&customer.Password,
		&customer.CusAddress,
		&customer.CusPh,
	); err != nil {
		// Handle authentication failure
		fmt.Println("Authentication failed")
		http.Redirect(w, r, "log.html", http.StatusSeeOther)
		return
	}

	// Authentication successful
	fmt.Println("Authentication successful")
	http.Redirect(w, r, "home.html", http.StatusSeeOther)
}
