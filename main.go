package main

import (
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	user "github.com/atomedgesoft/calendariq/controller/user"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/user", user.InsertUser).Methods("POST")
	router.HandleFunc("/user", user.GetUser).Methods("GET")
	router.HandleFunc("/user", user.UpdateUser).Methods(("PATCH"))
	router.HandleFunc("/user/{id}", user.GetUserById).Methods(("GET"))
	router.HandleFunc("/user/{id}", user.DeleteUserById).Methods("DELETE")
	c := cors.AllowAll()
	handler := c.Handler(router)
	log.Println("Server started on :80")
	log.Fatal(http.ListenAndServe(":80", handler))
}
