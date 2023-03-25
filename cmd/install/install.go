package install

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
)

var (
	tomasterDir = "/.tomaster"
	tomasterBin = "tomaster"
	symlinkName = "/usr/local/bin/tomaster"
)

func Run(home, binaryPath string, fs billy.Filesystem) error {
	targetDir := home + tomasterDir

	if s, _ := fs.Readlink(symlinkName); s != "" {
		fmt.Println("Already instaled, performing a reinstall")
	}

	fmt.Printf("\nCreating dir: %s\n", targetDir)
	if err := fs.MkdirAll(targetDir, 0644); err != nil {
		return err
	}

	tomasterbinDir := targetDir + "/" + tomasterBin
	fmt.Printf("Moving binary to: %s/tomaster\n", targetDir)
	if err := fs.Rename(binaryPath, tomasterbinDir); err != nil {
		return err
	}

	fmt.Printf("Creating symbolic link: %s -> %s\n", tomasterbinDir, symlinkName)
	if err := fs.Symlink(tomasterbinDir, symlinkName); err != nil && err.Error() != "file already exists" {
		return err
	}

	fmt.Printf("Symlink created at `%s`\nInstallation succesfull!\nYou can now call `tomaster` on shell\n\n", symlinkName)
	return nil
}
