package commands

import "fmt"

func requiresDatabaseErr(cmd string) error {
	return fmt.Errorf("/%s command requires a database", cmd)
}

func unknownDirectiveErr(s string) error {
	return fmt.Errorf("unknown directive: %s", s)
}
