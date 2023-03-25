package install_test

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"

	"github.com/luinunesmeli/goscriba/cmd/install"
)

func TestRun(t *testing.T) {
	t.Run("Create a new install", func(t *testing.T) {
		fs := memfs.New()
		f, _ := fs.TempFile("/Downloads", "tomaster_bin")
		f.Write([]byte("some data"))
		f.Close()

		assert.NoError(t, install.Run("/Users/luinunes", f.Name(), fs))
	})

	t.Run("Update an install", func(t *testing.T) {
		fs := memfs.New()
		f, _ := fs.TempFile("/Downloads", "tomaster_bin")
		f.Write([]byte("some data"))
		f.Close()

		f2, _ := fs.TempFile("/Users/luinunes/.tomaster", "tomaster_bin")
		f2.Write([]byte("some data"))
		f2.Close()

		fs.Symlink("/Users/luinunes/.tomaster", "/usr/local/bin/tomaster")

		assert.NoError(t, install.Run("/Users/luinunes", f.Name(), fs))
	})
}
