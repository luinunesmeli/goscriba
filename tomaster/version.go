package tomaster

import (
	"fmt"
	"strconv"
	"strings"
)

func NextReleases(latestTag string) (string, string, string, error) {
	parts := strings.Split(latestTag, ".")

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", "", err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", "", err
	}

	major, err := strconv.Atoi(strings.TrimPrefix(parts[0], "v"))
	if err != nil {
		return "", "", "", err
	}

	versionFmt := "%d.%d.%d"
	return fmt.Sprintf(versionFmt, major+1, 0, 0),
		fmt.Sprintf(versionFmt, major, minor+1, 0),
		fmt.Sprintf(versionFmt, major, minor, patch+1),
		nil
}
