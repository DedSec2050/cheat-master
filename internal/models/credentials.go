package models

import (
	"time"
)

// User represents a registered user with credentials
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // Consider encryption in production
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CredentialsStore manages user credentials
type CredentialsStore struct {
	Users map[string]*User `json:"users"`
}

// Job represents a course watching job
type Job struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CourseSlug  string    `json:"course_slug"`
	Status      string    `json:"status"` // pending, running, completed, failed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Progress    int       `json:"progress"` // percentage 0-100
	Error       string    `json:"error,omitempty"`
}

// JobQueue manages job scheduling
type JobQueue struct {
	Jobs map[string]*Job `json:"jobs"`
}
