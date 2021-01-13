package opnsense

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func JSONFields(b interface{}) []string {
	fields := []string{}
	val := reflect.ValueOf(b)

	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
			}
		}

		fields = append(fields, fieldName)
	}

	return fields
}

// func MapToStruct(input interface{}, output interface{}) error {
// 	config := &mapstructure.DecoderConfig{
// 		WeaklyTypedInput: true,
// 		Result:           &output,
// 	}

// 	decoder, err := mapstructure.NewDecoder(config)
// 	if err != nil {
// 		return err
// 	}

// 	err = decoder.Decode(input)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

/*
This function will help you to convert your object from struct to
map[string]interface{} based on your JSON tag in your structs.
https://gist.github.com/bxcodec/c2a25cfc75f6b21a0492951706bc80b8
*/
func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	if item == nil {
		return res
	}

	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()

		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}

	return res
}

// func StructToMap(data interface{}) (map[string]interface{}, error) {
// 	dataBytes, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mapData := make(map[string]interface{})

// 	err = json.Unmarshal(dataBytes, &mapData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mapData, nil
// }

func MapToStruct(input map[string]interface{}, output interface{}) error {
	dataBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	fmt.Printf("MapToStruct: %s\n", string(dataBytes))

	err = json.Unmarshal(dataBytes, &output)
	if err != nil {
		return err
	}

	return nil
}
