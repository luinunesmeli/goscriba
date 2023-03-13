package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	tomasterDir = "/.tomaster"
	tomasterBin = "tomaster"
	symlinkName = "/usr/local/bin/tomaster"
)

func Run() error {
	targetDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	targetDir += tomasterDir

	res, err := filepath.EvalSymlinks(symlinkName)
	if res != "" {
		return errors.New("ToMaster already installed, installation aborted!\n")
	}
	if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
		return err
	}

	fmt.Printf("The installation process will create a symlink at `%s`\n", symlinkName)
	if _, err := os.ReadDir(targetDir); err != nil {
		if err := os.Mkdir(targetDir, 0777); err != nil {
			return err
		}
	}

	path, err := os.Executable()
	if err != nil {
		return err
	}

	tomasterBinDir := fmt.Sprintf("%s/%s", targetDir, tomasterBin)
	fmt.Printf("Moving binary to `%s`\n", tomasterBinDir)
	if err := os.Rename(path, tomasterBinDir); err != nil {
		return err
	}

	if err := os.Symlink(tomasterBinDir, symlinkName); err != nil {
		return err
	}
	fmt.Printf("Symlink created at `%s`\nInstallation succesfull!\nYou can now call `tomaster` on shell\n", symlinkName)

	return nil
}
