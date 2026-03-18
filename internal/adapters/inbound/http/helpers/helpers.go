package helpers

import (
	"encoding/json"
	"strings"
)

func MapToStruct(m map[string]any, target any) error {
	jsonBytes, _ := json.Marshal(m)
	return json.Unmarshal(jsonBytes, target)
}

func IsTrimmedEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
