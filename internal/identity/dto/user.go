package dto

import "github.com/google/uuid"

type User struct {
	Id         int64
	Identifier uuid.UUID
	FirstName  string
	LastName   string
}
