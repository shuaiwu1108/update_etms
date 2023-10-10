package util

import (
	"database/sql"
)

func CustomQuery(DB *sql.DB, sqlString string) []map[string]interface{} {
	rows, _ := DB.Query(sqlString)
	defer rows.Close()
	var res []map[string]interface{}
	columns, _ := rows.Columns()
	vals := make([]interface{}, len(columns))
	valsPtr := make([]interface{}, len(columns))
	for i := range vals {
		valsPtr[i] = &vals[i]
	}
	for rows.Next() {
		_ = rows.Scan(valsPtr...)
		r := make(map[string]interface{})
		for i, v := range columns {
			if va, ok := vals[i].([]byte); ok {
				r[v] = string(va)
			} else {
				r[v] = vals[i]
			}
		}
		res = append(res, r)
	}
	return res
}
