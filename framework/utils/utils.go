package utils

import "encoding/json"

func isJson(s string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}

	return nil
}
