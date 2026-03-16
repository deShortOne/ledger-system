package common

import "errors"

var (
	ErrUserIdentifierNotFound = errors.New("user of that identifier does not exist")
	ErrUserDetailsNotFilledIn = errors.New("first name and last name must both be entered")
)
