package json

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Struct2Json 结构体转为json
func Struct2Json(obj interface{}) (str string, err error) {
	bt, err := json.Marshal(obj)
	if err != nil {
		return str, errors.New(fmt.Sprintf("[Struct2Json] transform error: %v", err))
	}
	return string(bt), nil
}

// Json2Struct json转为结构体
func Json2Struct(str string, obj interface{}) error {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return errors.New(fmt.Sprintf("[Json2Struct] transform error: %v", err))
	}
	return nil
}

// JsonI2Struct json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) error {
	JsonStr := str.(string)
	err := Json2Struct(JsonStr, obj)
	if err != nil {
		return err
	}
	return nil
}
