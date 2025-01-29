package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// generate statment string (SELECT x1, ...[v] FROM <table> ?WHERE <cond>)
func GenStmt(v any, table, cond string) (string, error) {
	// type control
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("passed v is not a struct")
	}

	keys := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		keys = append(keys, field.Name)
	}

	stmt := fmt.Sprintf("SELECT %s FROM %s", strings.Join(keys, " "), table)

	if cond != "" {
		stmt += fmt.Sprintf(" WHERE %s", cond)
	}

	return stmt, nil
}
