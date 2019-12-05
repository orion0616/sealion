package util

import (
	"encoding/json"
	"net/http"
)

// DecodeResponse decodes httpResponse to map[string]interface{}
func DecodeResponse(resp *http.Response) (map[string]interface{}, error) {
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
