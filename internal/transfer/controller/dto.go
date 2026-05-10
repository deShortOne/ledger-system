package controller

type DepositMoneyRequest struct {
	ToAccountId string
	Amount      float64
}

type TransferMoneyRequest struct {
	FromAccountId string
	ToAccountId   string
	Amount        float64
}
