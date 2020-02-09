package msg

import (
	"database/sql"
	"io"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// EelPic returns a single eel image from the eelbot database.
func EelPic() (io.Reader, error) {
	eelPics, err := EelBomb(1)
	if err != nil {
		return nil, err
	}

	return eelPics[0], nil
}

// EelBomb returns the given number of eel images from the eelbot database.
func EelBomb(num uint64) ([]io.Reader, error) {
	var eelPics []io.Reader

	db, err := sql.Open("sqlite3", sqliteDBString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlQuery := "SELECT FullPath FROM \"EelPics\" ORDER BY RANDOM() LIMIT " + strconv.FormatUint(num, 10) + ";"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eelPicPath string
		rows.Scan(&eelPicPath)

		eelPic, err := os.Open(eelPicPath)
		if err == nil {
			eelPics = append(eelPics, eelPic)
		}
	}

	return eelPics, nil
}
