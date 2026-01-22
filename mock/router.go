package main

import (
	"log"
	"net/http"
)

func main() {
	err := RunServer()
	if err != nil {
		log.Fatal(err)
	}
}

func RunServer() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/lookup/states/", StatesGetMany)     // expect swagger api docs
	router.HandleFunc("GET /api/v1/lookup/states/{code}", StateByCode) // expect swagger api docs

	server := &http.Server{
		Addr:    ":8000", // PORT from env
		Handler: router,
	}

	return server.ListenAndServe()

}
