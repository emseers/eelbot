package msg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (interpreter *Interpreter) getTaunt() (file io.Reader, filename string, err error) {
	sqlQuery := "SELECT path FROM taunts ORDER BY RANDOM() LIMIT 1;"
	row := interpreter.db.QueryRow(sqlQuery)

	var path string
	if err = row.Scan(&path); err != nil {
		return
	}

	filename = filepath.Base(path)
	file, err = os.Open(path)
	return
}

func (interpreter *Interpreter) getSpecificTaunt(num uint64) (file io.Reader, filename string, err error) {
	sqlQuery := fmt.Sprintf("SELECT path FROM taunts WHERE id=%d;", num)
	row := interpreter.db.QueryRow(sqlQuery)

	var path string
	if err = row.Scan(&path); err != nil {
		return
	}

	filename = filepath.Base(path)
	file, err = os.Open(path)
	return
}
