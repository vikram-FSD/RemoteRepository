package main

import (
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	user "github.com/atomedgesoft/calendariq/controller"
	"github.com/gorilla/mux"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/user", user.InsertUser).Methods("POST")
	router.HandleFunc("/getuser", user.GetUser).Methods("GET")

	log.Println("Server started on :80")
	log.Fatal(http.ListenAndServe(":80", router))

}
