package main

import (
	"log"
	"net/http"

	employee "example.com/basics-practice/API/controller"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/employees", employee.GetAllOrgDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", r))

}
