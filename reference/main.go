package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pacedotdev/oto/otohttp"
)

func main() {
	if err := RunServer(); err != nil {
		log.Fatal(err)
	}
}

func RunServer() error {
	router := otohttp.NewRouter()

	router.HandleTyped(
		"GET /api/v1/lookup/states/",
		otohttp.Wrap(http.HandlerFunc(StatesGetMany), nil, StatesResponse{}, otohttp.HandlerMeta{
			Service: "States",
			Method:  "GetMany",
			Summary: "List all states",
			Tags:    []string{"states"},
		}),
	)

	router.HandleTyped(
		"GET /api/v1/lookup/states/{code}",
		otohttp.Wrap(http.HandlerFunc(StateByCode), nil, StateResponse{}, otohttp.HandlerMeta{
			Service: "States",
			Method:  "GetByCode",
			Summary: "Get state by code",
			Tags:    []string{"states"},
		}),
	)

	if err := writeOpenAPI(router, "openapi.json"); err != nil {
		return err
	}
	if err := writeClient(router, "client.gen.js"); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /docs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs.html")
	})
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.json")
	})
	mux.HandleFunc("GET /client.gen.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client.gen.js")
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	fmt.Println("Listening on :8000")
	return server.ListenAndServe()
}

func writeOpenAPI(router *otohttp.Router, path string) error {
	data, err := router.OpenAPI()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func writeClient(router *otohttp.Router, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return router.WriteClientJS(f)
}
