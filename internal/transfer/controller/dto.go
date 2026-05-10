package controller

type TransferMoneyRequest struct {
	FromAccountId string
	ToAccountId   string
	Amount        float64
}
