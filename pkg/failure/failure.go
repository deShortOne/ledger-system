package failure

type Failure struct {
	error
	code             FailureCode
	failureType      FailureType
	technicalMessage string
}

func NewFailure(
	code FailureCode,
	failureType FailureType,
	userMessage error,
	technicalMessage string,
) *Failure {
	return &Failure{
		code:             code,
		failureType:      failureType,
		error:            userMessage,
		technicalMessage: technicalMessage,
	}
}

func (f Failure) GetCode() FailureCode {
	return f.code
}

func (f Failure) GetFailureType() FailureType {
	return f.failureType
}

func (f Failure) GetTechnicalMessage() string {
	return f.technicalMessage
}

func (f *Failure) Error() string {
	if f == nil {
		return ""
	}
	if f.error != nil {
		return f.error.Error()
	}
	return f.technicalMessage
}

func (f *Failure) Unwrap() error {
	if f == nil {
		return nil
	}
	return f.error
}
