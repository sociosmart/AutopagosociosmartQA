package utils

import "encoding/json"

func BoolAddr(b bool) *bool {
	return &b
}

func IntAddr(i int) *int {
	return &i
}

func StringAddr(s string) *string {
	return &s
}

func Transform[T any](from any) *T {
	data, _ := json.Marshal(from)

	decoded := new(T)

	json.Unmarshal(data, decoded)

	return decoded
}
