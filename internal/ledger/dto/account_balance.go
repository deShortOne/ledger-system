package dto

import "time"

type AccountBalance struct {
	AccountId        int64
	Availablebalance float64
	UpdatedAt        time.Time
}
