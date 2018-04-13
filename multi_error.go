package query

import (
	"fmt"
)

// MultiError 存储了多个编码或解码错误
type MultiError map[string]error

func (e MultiError) Error() string {
	s := ""
	for _, err := range e {
		s = err.Error()
		break
	}

	switch len(e) {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}

	return fmt.Sprintf("%s (and %d other errors)", s, len(e)-1)
}
