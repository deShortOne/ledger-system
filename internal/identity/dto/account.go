package dto

import (
	"github.com/google/uuid"
)

type Account struct {
	Id          int64
	Identifier  uuid.UUID
	UserId      int64
	CreatedAt   CustomTime
	AccountType string
	Currency    string
	Status      string
}
