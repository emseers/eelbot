package msg

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// getTaunt returns a single taunt given the taunt ID.
func (interpreter *Interpreter) getTaunt(tauntID int) (file io.Reader, filename string, err error) {
	files, err := ioutil.ReadDir(interpreter.tauntsFolder)
	if err != nil {
		return
	}

	tauntID--
	if tauntID < 0 {
		return nil, "", errors.New("taunt number too small")
	}
	if len(files) <= tauntID {
		return nil, "", errors.New("taunt number too big")
	}

	filename = files[tauntID].Name()
	file, err = os.Open(path.Join(interpreter.tauntsFolder, filename))
	return
}
