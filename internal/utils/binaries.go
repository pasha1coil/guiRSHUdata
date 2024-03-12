package utils

import "encoding/json"

func Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
