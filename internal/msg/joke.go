package msg

import (
	"database/sql"
	"strconv"
)

func (interpreter *Interpreter) getJoke() (line1 string, line2 sql.NullString, err error) {
	sqlQuery := "SELECT JokeText, JokeTextLine2 FROM \"EelJokes\" ORDER BY RANDOM() LIMIT 1;"
	row := interpreter.sqliteDB.QueryRow(sqlQuery)

	err = row.Scan(&line1, &line2)
	return
}

func (interpreter *Interpreter) getSpecificJoke(jokeNum uint64) (line1 string, line2 sql.NullString, err error) {
	sqlQuery := "SELECT JokeText, JokeTextLine2 FROM \"EelJokes\" WHERE JokeID=" + strconv.FormatUint(jokeNum, 10) + ";"
	row := interpreter.sqliteDB.QueryRow(sqlQuery)

	err = row.Scan(&line1, &line2)
	return
}
