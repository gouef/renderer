package handlers

import (
	"bytes"
	"html/template"
	"reflect"
	"strings"
)

func preprocessTemplate(tmpl string, data interface{}) (string, error) {
	var buffer bytes.Buffer
	lines := strings.Split(tmpl, "\n")
	for _, line := range lines {
		if strings.Contains(line, "n:range") {
			startIndex := strings.Index(line, "n:range")
			endIndex := strings.Index(line[startIndex:], `"`) + startIndex + 1
			rangeExpression := line[startIndex+9 : endIndex-1] // Např. ".Users"

			items, err := getFieldByName(data, rangeExpression)
			if err != nil {
				return "", err
			}

			v := reflect.ValueOf(items)
			if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
				for i := 0; i < v.Len(); i++ {
					item := v.Index(i).Interface()
					replacement := strings.ReplaceAll(line, `n:range="`+rangeExpression+`"`, "")
					replacement = strings.ReplaceAll(replacement, "{{.}}", template.HTMLEscapeString(item.(string)))
					buffer.WriteString(replacement + "\n")
				}
			}
		} else {
			buffer.WriteString(line + "\n")
		}
	}
	return buffer.String(), nil
}
func getFieldByName(data interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if field[0] == '.' {
		field = field[1:]
	}
	return v.FieldByName(field).Interface(), nil
}
