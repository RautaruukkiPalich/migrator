package migrator

import (
	"strings"
)

func validateTableName(table string) error {
	table = strings.ToLower(table)

	if len(strings.Split(table, " ")) != 1 ||
		strings.Contains(table, "drop ") ||
		strings.Contains(table, " ") ||
		strings.Contains(table, ";") ||
		strings.Contains(table, ",") ||
		strings.Contains(table, ".") {
			return ErrInvalidTablename
	}
	return nil
}

func validateSelectQuery(query string) error {
	query = strings.ToLower(query)

	if strings.Contains(query, "drop table") ||
		strings.Contains(query, "insert into") ||
		strings.Contains(query, "delete from") ||
		!strings.Contains(query, "select") {
			return ErrInvalidQuery
		}
	return nil
}