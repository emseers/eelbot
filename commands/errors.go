package commands

import "fmt"

func unknownDirectiveErr(s string) error {
	return fmt.Errorf("unknown directive: %s", s)
}
