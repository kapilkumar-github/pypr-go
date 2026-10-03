package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetContextValue(c *gin.Context, key string) string {
	value, exists := c.Get(key)
	if !exists {
		return ""
	}

	result, ok := value.(string)
	if !ok {
		return ""
	}

	return result
}

func GeneratePlaceholders(rowCount, columnCount int) string {
	var b strings.Builder

	param := 1

	for row := 0; row < rowCount; row++ {
		if row > 0 {
			b.WriteString(", ")
		}

		b.WriteString("(")

		for col := 0; col < columnCount; col++ {
			if col > 0 {
				b.WriteString(", ")
			}

			fmt.Fprintf(&b, "$%d", param)
			param++
		}

		b.WriteString(")")
	}

	return b.String()
}

func ThisOrDefault[T comparable](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

func SliceOrEmpty[T any](value *[]T) []T {
	if value == nil {
		return []T{}
	}
	return *value
}

var variableRegex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}\}`)

func ExtractVariableCodes(subject, body string) []string {
	seen := make(map[string]struct{})
	codes := make([]string, 0)

	extract := func(text string) {
		matches := variableRegex.FindAllStringSubmatch(text, -1)

		for _, match := range matches {
			code := match[1]

			if _, exists := seen[code]; exists {
				continue
			}

			seen[code] = struct{}{}
			codes = append(codes, code)
		}
	}

	extract(subject)
	extract(body)

	return codes
}
