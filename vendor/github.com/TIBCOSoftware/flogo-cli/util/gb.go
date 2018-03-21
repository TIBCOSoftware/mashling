package fgutil

import (
	"os"
	"os/exec"
	"path"
)

func IsGbProject(basePath string) bool {

	sourceDir := path.Join(basePath, "src")
	vendorDir := path.Join(basePath, "vendor", "src")

	info, err := os.Stat(sourceDir)

	if err != nil || !info.IsDir() {
		return false
	}

	info, err = os.Stat(vendorDir)

	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

//Env checker?
//IsProject(Path string) bool

// Gb structure that contains gb project paths
type Gb struct {
	BinPath        string
	SourcePath     string
	VendorPath     string
	CodeSourcePath string
}

// NewGb creates a new Gb struct
func NewGb(codePath string) *Gb {

	env := &Gb{}
	env.BinPath = "bin"
	env.SourcePath = "src"
	env.VendorPath = path.Join("vendor", "src")
	env.CodeSourcePath = path.Join("src", codePath)

	return env
}

// Init creates directories for the gb project
func (e *Gb) Init(createBin bool) {
	os.MkdirAll(e.SourcePath, 0777)
	os.MkdirAll(e.VendorPath, 0777)
	os.MkdirAll(e.CodeSourcePath, 0777)

	if createBin {
		os.MkdirAll(e.BinPath, 0777)
	}
}

// Installed indicates if gb is installed
func (e *Gb) Installed() bool {
	return ExecutableExists("gb")
}

// NewBinFilePath creates a new file Path in the bin directory
func (e *Gb) NewBinFilePath(fileName string) string {
	return path.Join(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Gb) VendorFetch(path string, version string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", path)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Gb) VendorFetchWithBranch(path string, version string, branch string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", "-branch", branch, path)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, "-branch", branch, path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorDeleteSilent performs a 'gb vendor delete' silently
func (e *Gb) VendorDeleteSilent(path string) error {
	cmd := exec.Command("gb", "vendor", "delete", path)

	return cmd.Run()
}

// VendorDelete performs a 'gb vendor delete'
func (e *Gb) VendorDelete(path string) error {
	cmd := exec.Command("gb", "vendor", "delete", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Build performs a 'gb build'
func (e *Gb) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
