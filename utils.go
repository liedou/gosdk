package gosdk

import "strings"

func ConvertEscapedChars(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(s, "\\", "\\\\"),
					"\n", "\\n"),
				"\t", "\\t"),
			"\r", "\\r"),
		"\"", "\\\"")
}
