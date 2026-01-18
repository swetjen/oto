package main

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
