package msg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func (interpreter *Interpreter) getImage() (file io.Reader, filename string, err error) {
	sqlQuery := "SELECT path FROM images ORDER BY RANDOM() LIMIT 1;"
	row := interpreter.db.QueryRow(sqlQuery)

	var path string
	if err = row.Scan(&path); err != nil {
		return
	}

	filename = filepath.Base(path)
	file, err = os.Open(path)
	return
}

func (interpreter *Interpreter) getSpecificImage(num uint64) (file io.Reader, filename string, err error) {
	sqlQuery := fmt.Sprintf("SELECT path FROM images WHERE id=%d;", num)
	row := interpreter.db.QueryRow(sqlQuery)

	var path string
	if err = row.Scan(&path); err != nil {
		return
	}

	filename = filepath.Base(path)
	file, err = os.Open(path)
	return
}
