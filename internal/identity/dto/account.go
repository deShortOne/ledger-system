package dto

import (
	"github.com/google/uuid"
)

type Account struct {
	Identifier     uuid.UUID
	UserIdentifier uuid.UUID
	CreatedAt      CustomTime
	AccountType    string
	Currency       string
	Status         string
}
