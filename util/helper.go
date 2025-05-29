package util

import (
	"fmt"
	"reflect"
	"strings"
)

func FormatWithTags(v any, fieldName string) string {
	var sb strings.Builder

	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("label")

		if strings.EqualFold(fieldType.Name, fieldName) {
			sb.WriteString(fmt.Sprintf("%s:%s", tag, field.String()))
		}
	}

	return sb.String()
}
