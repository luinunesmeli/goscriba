package generatetemplate

import (
	"fmt"

	"github.com/go-git/go-billy/v5"

	"github.com/luinunesmeli/goscriba/tomaster"
)

const projectTemplate = `---
changelog:
    path: changelog.md
    release_label: release
---
%s`

func Run(fs billy.Filesystem) error {
	f, err := fs.Create("./.tomaster")
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(fmt.Sprintf(projectTemplate, tomaster.ChangelogTemplate)))
	return err
}
