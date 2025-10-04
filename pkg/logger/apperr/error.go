package apperr

import "fmt"

func ErrorWrapper(place string, err error) error  {
	return fmt.Errorf("%s: %w", place, err)
}