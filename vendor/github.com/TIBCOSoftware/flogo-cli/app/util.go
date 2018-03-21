package app

import (
	"bufio"
	"fmt"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	CommandRegistry = cli.NewCommandRegistry()
)

func SetupNewProjectEnv() env.Project {
	return env.NewFlogoProject()
}

func SetupExistingProjectEnv(rootDir string) env.Project {

	proj := env.NewFlogoProject()

	if err := proj.Init(rootDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing flogo app project: %s\n\n", err.Error())
		os.Exit(2)
	}

	if err := proj.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening flogo app project: %s\n\n", err.Error())
		os.Exit(2)
	}

	return proj
}

func splitVersion(t string) (path string, version string) {

	idx := strings.LastIndex(t, "@")

	version = ""
	path = t

	if idx > -1 {
		v := t[idx+1:]

		if isValidVersion(v) {
			version = v
			path = t[0:idx]
		}
	}

	return path, version
}

//todo validate that "s" a valid semver
func isValidVersion(s string) bool {

	if s == "" {
		//assume latest version
		return true
	}

	if s[0] == 'v' && len(s) > 1 && isNumeric(string(s[1])) {
		return true
	}

	if isNumeric(string(s[0])) {
		return true
	}

	return false
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func Usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func cmdUsage(command cli.Command) {
	cli.CmdUsage("", command)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	options := CommandRegistry.CommandOptionInfos()
	options = append(options, cli.GetToolOptionInfos()...)

	fgutil.RenderTemplate(bw, usageTpl, options)
	bw.Flush()
}

var usageTpl = `Usage:

    flogo <command> [arguments]

Commands:
{{range .}}
    {{.Name | printf "%-12s"}} {{.Short}}{{end}}
`

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
