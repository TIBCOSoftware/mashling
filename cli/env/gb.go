package env

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"path"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

// ErrorNoManifest is generated when there is no manifest
var ErrorNoManifest = errors.New("no manifest")

type GbProject struct {
	BinDir         string
	RootDir        string
	SourceDir      string
	VendorDir      string
	VendorSrcDir   string
	CodeSourcePath string
	DockerBuild    bool
}

func NewGbProjectEnv() Project {

	env := &GbProject{}
	env.SourceDir = "src"
	env.VendorDir = "vendor"
	env.VendorSrcDir = path.Join("vendor", "src")

	return env
}

func (e *GbProject) Init(basePath string) error {

	exists := fgutil.ExecutableExists("gb")

	if !exists {
		return errors.New("gb not installed")
	}

	e.RootDir = basePath
	e.SourceDir = path.Join(basePath, "src")
	e.VendorDir = path.Join(basePath, "vendor")
	e.VendorSrcDir = path.Join(basePath, "vendor", "src")

	return nil
}

// Init creates directories for the gb project
func (e *GbProject) Create(createBin bool, vendorDir string) error {

	if _, err := os.Stat(e.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
	}

	os.MkdirAll(e.RootDir, os.ModePerm)
	os.MkdirAll(e.SourceDir, os.ModePerm)

	if vendorDir != "" {

		if strings.HasSuffix(vendorDir, string(os.PathSeparator)) {
			vendorDir = vendorDir[:len(vendorDir)-1]
		}

		if _, err := os.Stat(vendorDir); err != nil {
			return fmt.Errorf("Vendor directory '%s' not found", vendorDir)
		}

		if strings.HasSuffix(vendorDir, "vendor") {

			isGBVendor := false

			if _, err := os.Stat(path.Join(vendorDir, "src")); err == nil {
				if _, err := os.Stat(path.Join(vendorDir, "manifest")); err == nil {
					isGBVendor = true
				}
			}

			if isGBVendor {
				// is gb vendor dir, so copy to vendor dir
				fgutil.CopyDir(vendorDir, e.VendorDir)

			} else {
				// go vendor dir, so copy to vendor src sir
				fgutil.CopyDir(vendorDir, e.VendorSrcDir)
			}

		} else {
			fgutil.CopyDir(vendorDir, e.VendorSrcDir)
		}
	} else {
		os.MkdirAll(e.VendorSrcDir, os.ModePerm)
	}

	if createBin {
		e.BinDir = path.Join(e.RootDir, "bin")
		os.MkdirAll(e.BinDir, os.ModePerm)
	}

	return nil
}

// Open the project directory and validate its structure
func (e *GbProject) Open() error {

	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

	info, err = os.Stat(e.VendorSrcDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, vendor directory doesn't exists")
	}

	binDir := path.Join(e.RootDir, "bin")
	info, err = os.Stat(binDir)

	if err != nil || info.IsDir() {
		e.BinDir = binDir
	}

	return nil
}

func (e *GbProject) SetDockerBuild() {
	e.DockerBuild = true
}

func (e *GbProject) GetDockerBuild() bool {
	return e.DockerBuild
}

func (e *GbProject) GetBinDir() string {
	return e.BinDir
}

func (e *GbProject) GetRootDir() string {
	return e.RootDir
}

func (e *GbProject) GetSourceDir() string {
	return e.SourceDir
}

func (e *GbProject) GetVendorDir() string {
	return e.VendorDir
}

func (e *GbProject) GetVendorSrcDir() string {
	return e.VendorSrcDir
}

func (e *GbProject) InstallDependency(depPath string, version string) error {
	var cmd *exec.Cmd
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err == nil {
		//todo ignore installed dependencies for now
		//exists, return
		return nil
	}

	if version == "" {
		//if strings.HasPrefix(depPath,"github.com/TIBCOSoftware/flogo-") {
		//	cmd = exec.Command("gb", "vendor", "fetch", "-branch", "entrypoint", depPath)
		//} else {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
		//}
	} else {
		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	os.Chdir(e.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

//Restores dependecies using a manifest in the current working directory.
func (e *GbProject) RestoreDependency(manifest io.Reader) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//if manifest doesn't exists, no need to restore. return with error
	if manifest == nil {
		return ErrorNoManifest
	}

	var manifestcopy = path.Join(e.VendorDir, "manifest")

	// Create a manifest copy to the gb vendor dir
	mfTrg, err := os.Create(manifestcopy)
	defer mfTrg.Close()
	if err != nil {
		return err
	}

	// Copy the bytes to destination from source
	_, err = io.Copy(mfTrg, manifest)
	if err != nil {
		return err
	}

	// Flushes to disk
	err = mfTrg.Sync()
	if err != nil {
		return err
	}

	cmd = exec.Command("gb", "vendor", "restore")

	os.Chdir(e.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *GbProject) UninstallDependency(depPath string) error {

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err != nil {
		//todo ignore dependencies that are not installed for now
		//exists, return
		return nil
	}

	os.Chdir(e.RootDir)

	cmd := exec.Command("gb", "vendor", "delete", depPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *GbProject) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	os.Chdir(e.RootDir)
	fmt.Println(e.RootDir)

	return cmd.Run()
}

func IsGbProject(projectPath string) bool {

	sourceDir := path.Join(projectPath, "src")
	vendorDir := path.Join(projectPath, "vendor", "src")

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
//IsProject(path.Join string) bool

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
	return fgutil.ExecutableExists("gb")
}

// NewBinFilepath.Join creates a new file path.Join in the bin directory
func (e *Gb) NewBinFilePath(fileName string) string {
	return path.Join(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Gb) VendorFetch(depPath string, version string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorDeleteSilent performs a 'gb vendor delete' silently
func (e *Gb) VendorDeleteSilent(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)

	return cmd.Run()
}

// VendorDelete performs a 'gb vendor delete'
func (e *Gb) VendorDelete(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)
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

func (e *GbProject) GetAppDir() string {
	return ""
}
