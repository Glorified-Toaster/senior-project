package helpers

import "encoding/json"

func MapToStruct(m map[string]any, target any) error {
	jsonBytes, _ := json.Marshal(m)
	return json.Unmarshal(jsonBytes, target)
}
