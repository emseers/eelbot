package msg

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Joke returns a single joke from the eelbot database.
// If the joke is not multiline, the second string is left blank.
func Joke() (string, string, error) {
	db, err := sql.Open("sqlite3", sqliteDBString)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	sqlQuery := "SELECT JokeText, JokeTextLine2 FROM \"EelJokes\" ORDER BY RANDOM() LIMIT 1;"
	row := db.QueryRow(sqlQuery)

	var jokeLine1 string
	var jokeLine2 string
	row.Scan(&jokeLine1, &jokeLine2)
	return jokeLine1, jokeLine2, nil
}

// JokeSpecific returns a single joke given the joke ID in the eelbot database.
func JokeSpecific(jokeNum uint64) (string, string, error) {
	db, err := sql.Open("sqlite3", sqliteDBString)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	sqlQuery := "SELECT JokeText, JokeTextLine2 FROM \"EelJokes\" WHERE JokeID=" + strconv.FormatUint(jokeNum, 10) + ";"
	row := db.QueryRow(sqlQuery)

	var jokeLine1 string
	var jokeLine2 string
	row.Scan(&jokeLine1, &jokeLine2)
	return jokeLine1, jokeLine2, nil
}
