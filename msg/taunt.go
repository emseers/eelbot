package msg

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var tauntsFolder = "taunts"

func PlayTaunt(taunt int) (io.Reader, string, error) {
	files, err := ioutil.ReadDir(tauntsFolder)
	if err != nil {
		return nil, "", err
	} else if len(files) <= taunt {
		return nil, "", errors.New("Taunt number too big.")
	}

	file, err := os.Open(tauntsFolder + "/" + files[taunt].Name())
	if err != nil {
		return nil, "", err
	}

	return bufio.NewReader(file), files[taunt].Name(), nil
}
