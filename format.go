package log

import (
	"errors"
	"strings"
)

// Format is the logging format
type Format uint32

var (
	// FormatJSON tells the logger to write json structured messages
	FormatJSON Format = 1

	// FormatText tells the logger to write text structured key=value pairs messages
	FormatText Format = 2
)

func (l Format) String() (level string) {
	switch l {
	case FormatJSON:
		return "json"
	case FormatText:
		return "text"
	default:
		return "unknown"
	}
}

// ParseFormat parses the log format from a string
func ParseFormat(fmt string) (f Format, err error) {
	fmt = strings.ToLower(fmt)

	switch fmt {
	case "json":
		return FormatJSON, nil
	case "text":
		return FormatText, nil
	default:
		return Format(0), errors.New("unknown log format")
	}
}
