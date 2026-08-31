package common

import (
	"encoding/csv"
	"strings"
)

func ParseCSVLine(line string) []string {
	if hasLegacyUnquotedJSONObject(line) {
		return parseLegacyCSVLine(line)
	}

	reader := csv.NewReader(strings.NewReader(line))
	reader.FieldsPerRecord = -1
	fields, err := reader.Read()
	if err != nil {
		return parseLegacyCSVLine(line)
	}
	return fields
}

func hasLegacyUnquotedJSONObject(line string) bool {
	inQuotes := false
	atFieldStart := true

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				i++
				continue
			}
			if inQuotes {
				inQuotes = false
			} else if atFieldStart {
				inQuotes = true
			}
			atFieldStart = false
		case ',':
			if !inQuotes {
				atFieldStart = true
			}
		case '{':
			if !inQuotes && atFieldStart {
				return true
			}
			atFieldStart = false
		case ' ', '\t':
			// Keep looking for the first non-space character in the field.
		default:
			atFieldStart = false
		}
	}

	return false
}

func parseLegacyCSVLine(line string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false
	braceCount := 0

	for i := 0; i < len(line); i++ {
		char := line[i]
		if inQuotes {
			if char != '"' {
				current.WriteByte(char)
				continue
			}
			if i+1 < len(line) && line[i+1] == '"' {
				current.WriteByte('"')
				i++
				continue
			}
			inQuotes = false
			continue
		}

		switch char {
		case '"':
			if current.Len() == 0 {
				inQuotes = true
			} else {
				current.WriteByte(char)
			}
		case '{':
			braceCount++
			current.WriteByte(char)
		case '}':
			current.WriteByte(char)
			braceCount--
		case ',':
			if braceCount > 0 {
				current.WriteByte(char)
			} else {
				fields = append(fields, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}

	fields = append(fields, current.String())

	return fields
}
