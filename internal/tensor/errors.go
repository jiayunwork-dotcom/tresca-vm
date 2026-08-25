package tensor

import "errors"

var errPrincipalLength = errors.New("principal stress list must have length 3")

func IsLengthError(err error) bool {
	return errors.Is(err, errPrincipalLength)
}
