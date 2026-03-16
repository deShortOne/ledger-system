package contracts

import "errors"

var (
	ErrDoubleEntryViolated                  = errors.New("double entry has been violated")
	ErrOneOfTheAccountsDoNotHaveEnoughMoney = errors.New("one of the accounts does not have enough money")
)
