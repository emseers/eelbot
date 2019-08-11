package msg

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDBString = "EelbotDB.db"

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

func JokeSpecific(jokeNum string) (string, string, error) {
	db, err := sql.Open("sqlite3", sqliteDBString)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	sqlQuery := "SELECT JokeText, JokeTextLine2 FROM \"EelJokes\" WHERE JokeID=" + jokeNum + ";"
	row := db.QueryRow(sqlQuery)

	var jokeLine1 string
	var jokeLine2 string
	row.Scan(&jokeLine1, &jokeLine2)
	return jokeLine1, jokeLine2, nil
}
