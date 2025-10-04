package apperr

import "fmt"

func Wrap(opName string, err error, descr ...string) error {
	if len(descr) > 0 && descr[0] != "" {
			return fmt.Errorf("%s: %s: %w", opName, descr[0], err)
	}
	return fmt.Errorf("%s: %w", opName, err)
}