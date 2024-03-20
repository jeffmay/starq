package tsq

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: Use ["/pets"] instead of ./pets for keys with non-word characters.
func renderPath(path []string) string {
	if len(path) == 0 {
		return "."
	}
	res := new(strings.Builder)
	for i, p := range path {
		_, err := strconv.Atoi(p)
		if err == nil {
			if i == 0 {
				res.WriteRune('.')
			}
			res.WriteRune('[')
			res.WriteString(p)
			res.WriteRune(']')
		} else {
			res.WriteRune('.')
			res.WriteString(p)
		}
	}
	return res.String()
}

func unexpectedTypeError(path []string, err error, from AnyValue) error {
	return fmt.Errorf("%w at path %s in: %s", err, renderPath(append(from.PathFromTop(), path...)), from.Top().Pretty())
}
