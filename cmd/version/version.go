package version

import "fmt"

func Run(version string) error {
	fmt.Printf("Tomaster %s\n", version)
	return nil
}
