package util

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var UnsupportedFieldError = func(kind reflect.Kind) error {
	return errors.New(fmt.Sprintf("unsupported field type: %s", kind.String()))
}

func SerializeSection(v any) string {
	return SerializeFieldName(v, "")
}

func SerializeFieldName(v any, fieldName string) string {
	var sb strings.Builder

	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("label")

		if fieldName != "" && strings.EqualFold(fieldType.Name, fieldName) {
			str, err := formatFieldValue(field)
			if err != nil {
				log.Fatal(err)
			}
			sb.WriteString(fmt.Sprintf("%s:%s", tag, str))
			break
		} else if fieldName == "" {
			str, err := formatFieldValue(field)
			if err != nil {
				log.Fatal(err)
			}
			sb.WriteString(fmt.Sprintf("%s:%s", tag, str))
		}
	}

	return sb.String()
}

func formatFieldValue(field reflect.Value) (string, error) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		if field.Bool() {
			return "1", nil
		} else {
			return "0", nil
		}
	case reflect.Ptr:
		if field.IsNil() {
			return "", nil
		}
		elem := field.Elem()
		return formatFieldValue(elem)
	default:
		return "", UnsupportedFieldError(field.Kind())
	}
}
