package def

// GreeterService is a polite API for greeting people.
type GreeterService interface {
	// Greet prepares a lovely greeting.
	Greet(GreetRequest) GreetResponse

	// CreateUser registers a new user and returns the stored record.
	CreateUser(CreateUserRequest) CreateUserResponse
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
	// Auth is the bearer token for this request.
	Auth string `otoauth:"scheme=BearerAuth,in=header,name=Authorization,prefix=Bearer" json:"-"`
	// Name is the person to greet.
	// It is required.
	Name string
}

// GreetResponse is the response object containing a
// person's greeting.
type GreetResponse struct {
	// Greeting is a nice message welcoming somebody.
	Greeting string
}

// UUID represents a UUID encoded as a string.
type UUID string

// DateTime represents an RFC3339 timestamp.
type DateTime string

// UserType represents the user's role.
type UserType string

const (
	UserTypeAdmin  UserType = "admin"
	UserTypeMember UserType = "member"
	UserTypeGuest  UserType = "guest"
)

// Address contains mailing information for a user.
type Address struct {
	Line1      string
	Line2      string
	City       string
	Region     string
	PostalCode string
	Country    string
}

// User represents a person in the system.
type User struct {
	ID          UUID
	Name        string
	Email       string
	Type        UserType
	Address     Address
	CreatedAt   DateTime
	LastLoginAt *DateTime
	Tags        []string
	Metadata    map[string]string
}

// CreateUserRequest is the request object for GreeterService.CreateUser.
type CreateUserRequest struct {
	User        User
	InvitedBy   *UUID
	RequestedAt DateTime
}

// CreateUserResponse is the response object for GreeterService.CreateUser.
type CreateUserResponse struct {
	User           User
	WelcomeMessage string
}
