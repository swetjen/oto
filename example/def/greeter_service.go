package def

import "time"

// GreeterService is a polite API for greeting people.
type GreeterService interface {
	// Greet prepares a lovely greeting.
	Greet(GreetRequest) GreetResponse
	// GetUser retrieves a user by their ID.
	GetUser(GetUserRequest) GetUserResponse
	// CreateUser creates a new user account.
	CreateUser(CreateUserRequest) CreateUserResponse
	// ListUsers returns a paginated list of users.
	ListUsers(ListUsersRequest) ListUsersResponse
	// UpdateUserPreferences updates a user's preferences.
	UpdateUserPreferences(UpdateUserPreferencesRequest) UpdateUserPreferencesResponse
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
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

// GetUserRequest is the request for fetching a user.
type GetUserRequest struct {
	// UserID is the unique identifier for the user.
	// example: "550e8400-e29b-41d4-a716-446655440000"
	UserID string
}

// GetUserResponse contains the requested user.
type GetUserResponse struct {
	// User is the requested user details.
	User User
}

// CreateUserRequest contains the data for creating a new user.
type CreateUserRequest struct {
	// Email is the user's email address.
	// example: "john.doe@example.com"
	Email string
	// Username is the desired username.
	// example: "johndoe"
	Username string
	// FullName is the user's full name.
	// example: "John Doe"
	FullName string
	// Age is the user's age in years.
	// example: 30
	Age int
	// Profile contains optional profile information.
	Profile UserProfile
	// Tags are optional labels for the user.
	// example: ["premium", "early-adopter"]
	Tags []string
}

// CreateUserResponse contains the newly created user.
type CreateUserResponse struct {
	// User is the created user with assigned ID.
	User User
}

// ListUsersRequest contains pagination parameters.
type ListUsersRequest struct {
	// Page is the page number (1-indexed).
	// example: 1
	Page int
	// PageSize is the number of results per page.
	// example: 20
	PageSize int
	// SortBy is the field to sort by.
	// example: "created_at"
	SortBy string
	// SortDescending indicates descending sort order.
	// example: true
	SortDescending bool
	// Filter contains optional filter criteria.
	Filter UserFilter
}

// ListUsersResponse contains a page of users.
type ListUsersResponse struct {
	// Users is the list of users for this page.
	Users []User
	// TotalCount is the total number of users matching the filter.
	TotalCount int
	// HasMore indicates if there are more pages.
	HasMore bool
}

// UpdateUserPreferencesRequest contains preference updates.
type UpdateUserPreferencesRequest struct {
	// UserID is the user to update.
	// example: "550e8400-e29b-41d4-a716-446655440000"
	UserID string
	// Preferences contains the new preference values.
	Preferences UserPreferences
}

// UpdateUserPreferencesResponse confirms the update.
type UpdateUserPreferencesResponse struct {
	// Success indicates if the update was successful.
	Success bool
	// UpdatedAt is when the preferences were updated.
	UpdatedAt time.Time
}

// User represents a user in the system.
type User struct {
	// ID is the unique identifier (UUID format).
	// example: "550e8400-e29b-41d4-a716-446655440000"
	ID string
	// Email is the user's email address.
	// example: "john.doe@example.com"
	Email string
	// Username is the user's chosen username.
	// example: "johndoe"
	Username string
	// FullName is the user's display name.
	// example: "John Doe"
	FullName string
	// Age is the user's age in years.
	// example: 30
	Age int
	// Balance is the user's account balance.
	// example: 150.75
	Balance float64
	// IsActive indicates if the account is active.
	// example: true
	IsActive bool
	// Profile contains extended profile information.
	Profile UserProfile
	// Preferences contains user settings.
	Preferences UserPreferences
	// Tags are labels associated with the user.
	// example: ["premium", "verified"]
	Tags []string
	// CreatedAt is when the user was created.
	CreatedAt time.Time
	// UpdatedAt is when the user was last modified.
	UpdatedAt time.Time
	// LastLoginAt is when the user last logged in (optional).
	LastLoginAt *time.Time
}

// UserProfile contains extended profile information.
type UserProfile struct {
	// Bio is a short biography.
	// example: "Software developer and coffee enthusiast"
	Bio string
	// AvatarURL is the URL to the user's avatar image.
	// example: "https://example.com/avatars/johndoe.png"
	AvatarURL string
	// Location is the user's location.
	// example: "San Francisco, CA"
	Location string
	// Website is the user's personal website.
	// example: "https://johndoe.dev"
	Website string
	// SocialLinks contains links to social profiles.
	SocialLinks []SocialLink
}

// SocialLink represents a link to a social profile.
type SocialLink struct {
	// Platform is the social platform name.
	// example: "twitter"
	Platform string
	// URL is the profile URL.
	// example: "https://twitter.com/johndoe"
	URL string
}

// UserPreferences contains user settings.
type UserPreferences struct {
	// Theme is the UI theme preference.
	// example: "dark"
	Theme string
	// Language is the preferred language code.
	// example: "en-US"
	Language string
	// EmailNotifications enables email notifications.
	// example: true
	EmailNotifications bool
	// WeeklyDigest enables weekly summary emails.
	// example: false
	WeeklyDigest bool
	// Timezone is the user's timezone.
	// example: "America/Los_Angeles"
	Timezone string
}

// UserFilter contains criteria for filtering users.
type UserFilter struct {
	// IsActive filters by active status.
	IsActive *bool
	// MinAge filters users with age >= this value.
	MinAge *int
	// MaxAge filters users with age <= this value.
	MaxAge *int
	// HasTag filters users that have this tag.
	HasTag string
	// CreatedAfter filters users created after this time.
	CreatedAfter *time.Time
	// CreatedBefore filters users created before this time.
	CreatedBefore *time.Time
}
