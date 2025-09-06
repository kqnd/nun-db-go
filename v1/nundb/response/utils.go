package response

import (
	"strings"
)

func splitResponse(message string) []string {
	no_left_brackets := strings.ReplaceAll("[", "", message)
	no_brackets := strings.ReplaceAll("]", "", no_left_brackets)
	split := strings.Split(no_brackets, " ")
	return split
}
