package entity

import (
	"time"
)

type UserRole string

const (
	Employee  UserRole = "employee"
	Moderator UserRole = "moderator"
)

type User struct {
	ID               string    `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	HashedPassword   string    `json:"-" db:"hashed_password"`
	Role             UserRole  `json:"role" db:"role"`
	RegistrationDate time.Time `json:"-" db:"registration_date"`
}
