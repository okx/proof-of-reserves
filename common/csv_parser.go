package common

import (
	"strings"
)

func ParseCSVLine(line string) []string {
	var fields []string
	var current strings.Builder
	inJSON := false
	braceCount := 0

	for _, char := range line {
		switch char {
		case '{':
			braceCount++
			if braceCount == 1 {
				inJSON = true
			}
			current.WriteRune(char)
		case '}':
			current.WriteRune(char)
			braceCount--
			if braceCount == 0 {
				inJSON = false
			}
		case ',':
			if inJSON {
				current.WriteRune(char)
			} else {
				fields = append(fields, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	fields = append(fields, current.String())

	return fields
}
