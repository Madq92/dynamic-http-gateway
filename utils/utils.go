package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type KipError struct {
	Msg string
}

func (s KipError) Error() string {
	return s.Msg
}

func GetMap(entity interface{}) (map[string]interface{}, error) {
	switch value := entity.(type) {
	case map[string]interface{}:
		return value, nil
	default:
		return nil, KipError{Msg: fmt.Sprintf("%v to map error", entity)}
	}
}

func GetString(entity interface{}) (string, error) {
	switch value := entity.(type) {
	case string:
		return value, nil
	default:
		return "", KipError{Msg: fmt.Sprintf("%v to string error", entity)}
	}
}

func GetFloat64(entity interface{}) (float64, error) {
	switch value := entity.(type) {
	case float64:
		return value, nil
	default:
		return 0, KipError{Msg: fmt.Sprintf("%v to float64 error", entity)}
	}
}

func DecodeJson(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}
