package utils

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

/*
- Не обязательный параметры будут передаваться через указатель и не обязательно должны передаваться в parseParamValue
- Обязательные нужно заполнять в любом случае или выдавать ошибку

- Проверка на limit и page будет отдельной функцией validateListsParam

Для каждого поля:
Получить param тег и получить значение из url.Values
Если нет значение из url.Values =>
	- если обязательное поле, то ошибка
Если есть =>
	- распаршиваем значение, записываем в поле
*/

func ParseParams(input interface{}, params url.Values) error {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		return errors.New("error_parse_params")
	}

	inputStruct := val.Elem()

	for i := 0; i < inputStruct.NumField(); i++ {
		field := inputStruct.FieldByIndex([]int{i})

		paramTagValue := inputStruct.Type().Field(i).Tag.Get("param")
		if paramTagValue == "" {
			return errors.New("error_absent_param_tag")
		}

		if paramTagValue == "-" {
			// fmt.Println("ParseParams: paramTagValue == \"-\"")
			continue
		}

		// // Проверяем существует ли параметр в url.Values // TODO: сделать проверку на bool без значения
		// if _, exists := params[paramTagValue]; !exists {
		// 	if field.Type().Kind() != reflect.Ptr {
		// 		return errors.New("required_field_without_value")
		// 	}
		// 	continue
		// }

		values := params[paramTagValue]
		if len(values) == 0 || values[0] == "" {
			if field.Type().Kind() != reflect.Ptr {
				return errors.New("required_field_without_value")
			}
			continue
		}

		if values[0] != "" { // Потенциальная ошибка
			if err := parseParamValue(field, values[0]); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseParamValue(field reflect.Value, paramStringValue string) error {
	is_ptr := false
	t := field.Type().Kind()
	if t == reflect.Ptr {
		t = field.Type().Elem().Kind()
		is_ptr = true
	}

	switch t {
	case reflect.Int64:
		val, err := strconv.ParseInt(paramStringValue, 10, 64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Uint:
		val64, err := strconv.ParseUint(paramStringValue, 10, 64)
		val := uint(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Int:
		val64, err := strconv.ParseInt(paramStringValue, 10, 64)
		val := int(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Int8:
		val64, err := strconv.ParseInt(paramStringValue, 10, 8)
		val := int8(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Int16:
		val64, err := strconv.ParseInt(paramStringValue, 10, 16)
		val := int16(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Int32:
		val64, err := strconv.ParseInt(paramStringValue, 10, 32)
		val := int32(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Uint8:
		val64, err := strconv.ParseUint(paramStringValue, 10, 8)
		val := uint8(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Uint16:
		val64, err := strconv.ParseUint(paramStringValue, 10, 16)
		val := uint16(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Uint32:
		val64, err := strconv.ParseUint(paramStringValue, 10, 32)
		val := uint32(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Float64:
		val, err := strconv.ParseFloat(paramStringValue, 64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Float32:
		val64, err := strconv.ParseFloat(paramStringValue, 32)
		val := float32(val64)
		if err != nil {
			return errors.New("error_parse_params")
		}

		if is_ptr {
			field.Set(reflect.ValueOf(&val))
		} else {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.String:
		if is_ptr {
			field.Set(reflect.ValueOf(&paramStringValue))
		} else {
			field.Set(reflect.ValueOf(paramStringValue))
		}
	case reflect.Slice: // Если срез, то срез строк
		strNums := strings.Split(paramStringValue, ",")
		if is_ptr {
			field.Set(reflect.ValueOf(&strNums))
		} else {
			field.Set(reflect.ValueOf(strNums))
		}
	case reflect.Bool:
		switch strings.ToLower(paramStringValue) {
		case "true":
			if is_ptr {
				boolVal := true
				field.Set(reflect.ValueOf(&boolVal))
			} else {
				field.Set(reflect.ValueOf(true))
			}
		case "false":
			if is_ptr {
				boolVal := false
				field.Set(reflect.ValueOf(&boolVal))
			} else {
				field.Set(reflect.ValueOf(false))
			}
		default:
			return fmt.Errorf("invalid boolean param: %s", paramStringValue)
		}
	default:
		// fmt.Println(field.Type().Kind())
		return errors.New("error_unknowing_param_type")
	}
	return nil
}
