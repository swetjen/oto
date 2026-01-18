package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (greeterService) GetUser(ctx context.Context, r GetUserRequest) (*GetUserResponse, error) {
	now := time.Now()
	return &GetUserResponse{
		User: User{
			ID:        r.UserID,
			Email:     "user@example.com",
			Username:  "testuser",
			FullName:  "Test User",
			Age:       30,
			Balance:   100.50,
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (greeterService) CreateUser(ctx context.Context, r CreateUserRequest) (*CreateUserResponse, error) {
	now := time.Now()
	return &CreateUserResponse{
		User: User{
			ID:        "generated-uuid-here",
			Email:     r.Email,
			Username:  r.Username,
			FullName:  r.FullName,
			Age:       r.Age,
			IsActive:  true,
			Profile:   r.Profile,
			Tags:      r.Tags,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (greeterService) ListUsers(ctx context.Context, r ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Users:      []User{},
		TotalCount: 0,
		HasMore:    false,
	}, nil
}

func (greeterService) UpdateUserPreferences(ctx context.Context, r UpdateUserPreferencesRequest) (*UpdateUserPreferencesResponse, error) {
	return &UpdateUserPreferencesResponse{
		Success:   true,
		UpdatedAt: time.Now(),
	}, nil
}

func main() {
	var greeterService greeterService
	server := otohttp.NewServer()
	RegisterGreeterService(server, greeterService)
	http.Handle("/oto/", server)
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
