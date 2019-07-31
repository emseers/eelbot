package msg

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Joke() (string, string, error) {
	db, err := sql.Open("sqlite3", "./EelBotDB.db")
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

	// jokes := [10]string{
	// 	"Can february march?\nNo, but april may",
	// 	"You can tune a piano, but you can't tuna fish!\nAmirite?",
	// 	"Did you hear about the kidnapping?\nHe woke up",
	// 	"See this eel?\nNow that's a moray",
	// 	"A ghost walks into a bar and orders a vodka and coke\nThe barman says \"we don't serve spirits\"",
	// 	"I often forget where i put my boomerang\nBut then it always comes back to me",
	// 	"Civil war jokes?\nI General Lee don't find them funny",
	// 	"Did you hear about the two guys who stole the calendar?\nThey both got 6 months each",
	// 	"Why is a mexican midget called a paragraph?\nBecause he's too short to be an essay",
	// 	"Why don't fish play basketball?\nBecause they're afraid of the net",
	// }

	// randJoke := strings.Split(jokes[rand.Intn(len(jokes))], "\n")
	// return randJoke[0], randJoke[1], nil
}
