package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	Name string
	Rows []Row
}

type Row struct {
	Column string
	Value  string
}

func selectInnerJoinWhereOrderLimit(selectTable Table, joinTable Table, selectField string, joinField string, where string, order string, limit int) (table Table) {
	querySelect := "SELECT "

	// Add all selected fields from selectTable
	for _, row := range selectTable.Rows {
		querySelect += selectTable.Name + "." + row.Column + " as " + selectTable.Name + "_" + row.Column + ", "
	}

	// Add all selected fields from joinTable
	for _, row := range joinTable.Rows {
		querySelect += joinTable.Name + "." + row.Column + " as " + selectTable.Name + "_" + row.Column + ", "
	}

	// Trim remain apostrophe
	querySelect = strings.TrimSuffix(querySelect, ", ")

	/*querySelect = `feeditemhash, feeds.FeedType, feeds.title as feedtitle,
	feed_items.title, Content, Categories,
	feed_items.description, feeds.link as feedlink, feed_items.link,
	feed_items.updated, feed_items.published,
	feeds.author as feedauthor, feed_items.author, language, copyright`*/

	querySelect += " FROM " + selectTable.Name
	querySelect += " INNER JOIN " + joinTable.Name + " ON " + selectTable.Name + "." + selectField + "=" + joinTable.Name + "." + joinField
	queryWhere := " WHERE " + where
	queryOrder := " ORDER BY " + order
	queryLimit := " LIMIT " + strconv.Itoa(limit)

	sqlQuery := querySelect + queryWhere + queryOrder + queryLimit + ";"

	fmt.Println(sqlQuery)

	// Perform query, return error
	rows, sqlErr := db.Query(sqlQuery)
	checkError(sqlErr)
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	tempRows := []Row{}
	//v := reflect.ValueOf(person)
	//f := v.FieldByName("address")

	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		checkError(err)

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}

			thisRow := Row{
				col,
				v.(string),
			}

			tempRows = append(tempRows, thisRow)
			//fmt.Println(col, v)
		}

		table.Rows = tempRows
		//results = append(results, tempTable)
	}
	err := rows.Err()
	checkError(err)

	return table
}

func insertUpdateUnique(table Table) (err error) {
	// Insert or update on unique already present after building statement.
	queryInsert := "INSERT IGNORE INTO " + table.Name + " ( "
	queryValues := "VALUES ( "
	queryUpdate := "ON DUPLICATE KEY UPDATE "

	// for each colum and row
	for _, row := range table.Rows {
		// Insert NULL on absent value
		if row.Value == "" {
			queryInsert += "`" + row.Column + "`, "
			queryValues += "NULL, "
			queryUpdate += "`" + row.Column + "` = NULL, "
		} else {
			queryInsert += "`" + row.Column + "`, "
			queryValues += "'" + MysqlRealEscapeString(row.Value) + "', "
			queryUpdate += "`" + row.Column + "` = '" + MysqlRealEscapeString(row.Value) + "', "
		}
	}

	// Trim remain apostrophe
	queryInsert = strings.TrimSuffix(queryInsert, ", ")
	queryValues = strings.TrimSuffix(queryValues, ", ")
	queryUpdate = strings.TrimSuffix(queryUpdate, ", ")

	// Add remaining parantheses and concat
	queryInsert += " ) "
	queryValues += " ) "
	sqlQuery := queryInsert + queryValues + queryUpdate + ";"

	// Perform query, return error
	_, sqlErr := db.Exec(sqlQuery)
	checkError(sqlErr)

	return err
}

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
