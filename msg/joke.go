package msg

import (
	"database/sql"
	"fmt"
)

func (interpreter *Interpreter) getJoke() (line1 string, line2 sql.NullString, err error) {
	sqlQuery := "SELECT text, punchline FROM jokes ORDER BY RANDOM() LIMIT 1;"
	row := interpreter.db.QueryRow(sqlQuery)

	err = row.Scan(&line1, &line2)
	return
}

func (interpreter *Interpreter) getSpecificJoke(num uint64) (line1 string, line2 sql.NullString, err error) {
	sqlQuery := fmt.Sprintf("SELECT text, punchline FROM jokes WHERE id=%d;", num)
	row := interpreter.db.QueryRow(sqlQuery)

	err = row.Scan(&line1, &line2)
	return
}
