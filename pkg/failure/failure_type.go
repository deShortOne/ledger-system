package failure

type FailureType int

const (
	GeneralFailure     FailureType = 1
	NotFound           FailureType = 2
	Validation         FailureType = 3
	AccessUnauthorised FailureType = 4
)
