package dto

import "github.com/google/uuid"

type User struct {
	Identifier uuid.UUID
	FirstName  string
	LastName   string
}
