package view

import (
	"errors"
	"strconv"

	"github.com/mritd/bubbles/prompt"
)

func newDaysInput() *prompt.Model {
	return &prompt.Model{
		Prompt: "Release days: ",
		ValidateFunc: func(value string) error {
			days, err := strconv.Atoi(value)
			if err != nil {
				return errors.New("invalid value, only integers are valid")
			}
			if days < 0 {
				return errors.New("negative days not permited")
			}
			return nil
		},
	}
}
