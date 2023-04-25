package tomaster_test

//import (
//	"io"
//	"testing"
//
//	"github.com/go-git/go-billy/v5/memfs"
//	"github.com/go-git/go-git/v5"
//	"github.com/go-git/go-git/v5/storage/memory"
//	"github.com/stretchr/testify/assert"
//
//	"github.com/luinunesmeli/goscriba/pkg/config"
//	"github.com/luinunesmeli/goscriba/tomaster"
//)
//
//var expected = `# Changelog
//
//## Version 1.2.3
//**Created at 2023-04-05 by **
//
//### Enhancements
//* [PR#9](https://github.com/owner/repo/pull/9) by author
//
//
//## Version 1.2.2
//**Created at 2023-04-04 by **
//
//### Enhancements
//* [PR#8](https://github.com/owner/repo/pull/8) by author
//`
//
//var fileContent = `# Changelog
//
//## Version 1.2.2
//**Created at 2023-04-04 by **
//
//### Enhancements
//* [PR#8](https://github.com/owner/repo/pull/8) by author
//`
//
//func TestChangelog_WriteChangelog(t *testing.T) {
//	fs := memfs.New()
//	fileCreate, _ := fs.Create("changelog.md")
//	defer fileCreate.Close()
//	fileCreate.Write([]byte(fileContent))
//
//	r, _ := git.Init(memory.NewStorage(), fs)
//	tree, _ := r.Worktree()
//
//	changelog := tomaster.NewChangelog(config.Config{
//		Changelog: "changelog.md",
//	}, tree)
//
//	res := changelog.WriteChangelog().Run(tomaster.Session{
//		ChosenVersion: "1.2.3",
//		PRs: []tomaster.PR{
//			{
//				PRType: tomaster.Enhancement,
//				Title:  "PR#9",
//				PRLink: "https://github.com/owner/repo/pull/9",
//				Author: "author",
//				Number: 95,
//			},
//		},
//	})
//
//	fileRead, _ := tree.Filesystem.Open("changelog.md")
//	b, _ := io.ReadAll(fileRead)
//
//	assert.NoError(t, res.Err)
//	assert.Equal(t, expected, string(b))
//}
