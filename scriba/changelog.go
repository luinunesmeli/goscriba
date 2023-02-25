package scriba

import (
	"errors"
	"os"
)

func changelogExist(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.New("changelog not found on path")
	}
	return nil
}
