package controller

type NewAccountRequest struct {
	AccountType string `json:"account_type" example:"checking,savings"`
	Currency    string `json:"currency" example:"GBP,USD,EUR"`
}

type NewUserRequst struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetAccountsForUserResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}

type AccountResponse struct {
	Identifier  string `json:"id"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created" example:"2024-02-29"`
}
