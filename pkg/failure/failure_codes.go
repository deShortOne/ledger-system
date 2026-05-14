package failure

type FailureCode int

const (
	// 0001 - 1000: General errors
	UnknownError           FailureCode = 1
	ConversionError        FailureCode = 2
	UnknownRepositoryError FailureCode = 3

	// 1001 - 2000: Identity
	UserValidationError     FailureCode = 1001
	UserNotFound            FailureCode = 1002
	IdentityRepositoryError FailureCode = 1003

	// 2001 - 3000: Ledger
	NotEnoughMoneyToDebit FailureCode = 2001
	DoubleEntryViolated   FailureCode = 2002
	LedgerRepositoryError FailureCode = 2003
	LedgerNotFound        FailureCode = 2004

	// 3001 - 4000: Platform
	PlatformRepositoryError FailureCode = 3001

	// 4001 - 5000: Transfer
	TransferRequestNotFound FailureCode = 4001
	TransferRepositoryError FailureCode = 4002
)
