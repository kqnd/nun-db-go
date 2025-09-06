package response

import (
	"strings"
)

func splitResponse(message string) []string {
	s := strings.TrimSpace(message)
	s = strings.Trim(s, "[]")
	return strings.Fields(s)
}
