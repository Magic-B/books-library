package op

import "strings"

// Operation name
func Namespace(namespace string) func(op ...string) string {
	return func(op ...string) string {
		if len(op) == 0 {
			return namespace
		}
		return namespace + "." + strings.Join(op, ".")
	}
}
