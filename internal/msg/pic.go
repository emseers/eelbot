package msg

import (
	"io"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func (interpreter *Interpreter) getEelPic() (file io.Reader, err error) {
	eelPics, err := interpreter.getEelBomb(1)
	if err != nil {
		return
	}

	file = eelPics[0]
	return
}

func (interpreter *Interpreter) getEelBomb(num uint64) (files []io.Reader, err error) {
	sqlQuery := "SELECT FullPath FROM \"EelPics\" ORDER BY RANDOM() LIMIT " + strconv.FormatUint(num, 10) + ";"
	rows, err := interpreter.sqliteDB.Query(sqlQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var eelPicPath string
		err1 := rows.Scan(&eelPicPath)
		if err1 != nil {
			err = err1
			return
		}

		eelPic, err2 := os.Open(eelPicPath)
		if err2 != nil {
			err = err2
			return
		}

		files = append(files, eelPic)
	}

	return
}
