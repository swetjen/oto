package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pacedotdev/oto/otohttp"
)

//go:generate ./generate.sh

// greeterService implements the generated GreeterService interface.
type greeterService struct{}

func (greeterService) Greet(ctx context.Context, r GreetRequest) (*GreetResponse, error) {
	resp := &GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s.", r.Name),
	}
	return resp, nil
}

func (greeterService) CreateUser(ctx context.Context, r CreateUserRequest) (*CreateUserResponse, error) {
	user := r.User
	if user.ID == "" {
		user.ID = UUID(fmt.Sprintf("user-%d", time.Now().UnixNano()))
	}
	if user.Type == "" {
		user.Type = UserTypeMember
	}
	if user.CreatedAt == "" {
		user.CreatedAt = DateTime(time.Now().UTC().Format(time.RFC3339Nano))
	}
	resp := &CreateUserResponse{
		User:           user,
		WelcomeMessage: fmt.Sprintf("Welcome, %s!", user.Name),
	}
	return resp, nil
}

const demoAuthToken = "demo-token"

type bearerGuard struct{}

func (bearerGuard) Authorize(ctx context.Context, r *http.Request, spec AuthSpec) error {
	header := r.Header.Get(spec.Name)
	if header == "" {
		return fmt.Errorf("missing auth token")
	}
	token := header
	if spec.Prefix != "" {
		prefix := spec.Prefix + " "
		if !strings.HasPrefix(header, prefix) {
			return fmt.Errorf("invalid auth token")
		}
		token = strings.TrimPrefix(header, prefix)
	}
	if token != demoAuthToken {
		return fmt.Errorf("invalid auth token")
	}
	return nil
}

func main() {
	var greeterService greeterService
	server := otohttp.NewServer()
	RegisterGreeterServiceWithAuth(server, greeterService, bearerGuard{})
	http.Handle("/oto/", server)
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	http.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs.html")
	})
	http.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./openapi.yaml")
	})
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// statusCodeHandler is useful for testing the server by returning a
// specific HTTP status code.
//  http.Handle("/", statusCodeHandler(http.StatusInternalServerError))
type statusCodeHandler int

func (c statusCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(c))
	io.WriteString(w, http.StatusText(int(c)))
}
