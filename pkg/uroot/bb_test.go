package uroot

import (
	"go/ast"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
)

func TestBBBuild(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := BuildOpts{
		Env: golang.Default(),
		Packages: []string{
			"github.com/u-root/u-root/pkg/uroot/test/foo",
			"github.com/u-root/u-root/cmds/rush",
		},
		TempDir: dir,
	}
	af := NewArchiveFiles()
	if err := BBBuild(af, opts); err != nil {
		t.Error(err)
	}

	var mustContain = []string{
		"init",
		"bbin/rush",
		"bbin/foo",
	}
	for _, name := range mustContain {
		if !af.Contains(name) {
			t.Errorf("expected files to include %q", name)
		}
	}

}

func findFile(filemap map[string]*ast.File, basename string) *ast.File {
	for name, f := range filemap {
		if filepath.Base(name) == basename {
			return f
		}
	}
	return nil
}

func TestPackageRewriteFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bin := filepath.Join(dir, "foo")
	if err := BuildBusybox(golang.Default(), []string{"github.com/u-root/u-root/pkg/uroot/test/foo"}, bin); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("foo failed: %v %v", string(o), err)
	}
}
