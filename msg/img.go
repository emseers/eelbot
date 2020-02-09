package msg

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func EelPic() (io.Reader, error) {
	db, err := sql.Open("sqlite3", sqliteDBString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlQuery := "SELECT FullPath FROM \"EelPics\" ORDER BY RANDOM() LIMIT 1;"
	row := db.QueryRow(sqlQuery)

	var eelPicPath string
	row.Scan(&eelPicPath)

	eelPic, err := os.Open(eelPicPath)
	if err != nil {
		return nil, err
	}

	return eelPic, nil
}
