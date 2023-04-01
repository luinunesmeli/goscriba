package uninstall

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
)

var (
	tomasterDir = "/.tomaster"
	symlinkName = "/usr/local/bin/tomaster"
)

func Run(home string, fs billy.Filesystem) error {
	targetDir := home + tomasterDir

	if s, _ := fs.Readlink(symlinkName); s != "" {
		fmt.Println("Removing symbolic link...")
		if err := fs.Remove(symlinkName); err != nil {
			return err
		}
	}

	fmt.Println("Removing tomaster files...")
	files, err := fs.ReadDir(targetDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if err = fs.Remove(targetDir + "/" + file.Name()); err != nil {
			return err
		}
	}

	if err = fs.Remove(targetDir); err != nil {
		return err
	}

	fmt.Println("ToMaster is no longer on your system!")

	return nil
}
