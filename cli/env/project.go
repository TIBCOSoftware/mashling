package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"io/ioutil"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

type MashlingProject struct {
	BinDir             string
	RootDir            string
	SourceDir          string
	VendorDir          string
	VendorSrcDir       string
	CodeSourcePath     string
	AppDir             string
	FileDescriptorPath string
	DockerBuild        bool
}

func NewMashlingProject() Project {
	return &MashlingProject{}
}

func (e *MashlingProject) SetDockerBuild() {
	e.DockerBuild = true
}

func (e *MashlingProject) GetDockerBuild() bool {
	return e.DockerBuild
}

func (e *MashlingProject) Init(rootDir string) error {

	exists := fgutil.ExecutableExists("dep")

	if !exists {
		return errors.New("dep not installed")
	}
	e.RootDir = rootDir
	e.SourceDir = filepath.Join(e.RootDir, "src")
	e.BinDir = filepath.Join(e.RootDir, "bin")
	return nil
}

// Create creates directories for the project
func (e *MashlingProject) Create(createBin bool, vendorDir string) error {

	if _, err := os.Stat(e.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
	}

	os.MkdirAll(e.RootDir, os.ModePerm)
	os.MkdirAll(e.SourceDir, os.ModePerm)

	return nil
}

// Open the project directory and validate its structure
func (e *MashlingProject) Open() error {

	// Check root dir
	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	// Check source dir
	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

	// Check file descriptor
	fd := filepath.Join(e.RootDir, config.FileDescriptor)
	_, err = os.Stat(fd)

	if err != nil {
		return fmt.Errorf("Invalid project, file descriptor '%s' doesn't exists", fd)
	}
	e.FileDescriptorPath = fd

	fdbytes, err := ioutil.ReadFile(fd)
	if err != nil {
		return fmt.Errorf("Invalid reading file descriptor '%s' error: %s", fd, err)
	}

	descriptor, err := ParseAppDescriptor(string(fdbytes))
	if err != nil {
		return fmt.Errorf("Invalid parsing file descriptor '%s' error: %s", fd, err)
	}

	appName := descriptor.Name

	// Validate that there is an app dir
	e.AppDir = filepath.Join(e.SourceDir, appName)
	info, err = os.Stat(e.AppDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid project, app directory '%s' doesn't exists", e.AppDir)
	}

	e.VendorDir = filepath.Join(e.AppDir, "vendor")
	e.VendorSrcDir = e.VendorDir

	return nil
}

func (e *MashlingProject) GetBinDir() string {
	return e.BinDir
}

func (e *MashlingProject) GetRootDir() string {
	return e.RootDir
}

func (e *MashlingProject) GetSourceDir() string {
	return e.SourceDir
}

func (e *MashlingProject) GetVendorDir() string {
	return e.VendorDir
}

func (e *MashlingProject) GetVendorSrcDir() string {
	return e.VendorSrcDir
}

// GetAppDir returns the directory of the app
func (e *MashlingProject) GetAppDir() string {
	return e.AppDir
}

// Deprecated, dependency managements responsibility
func (e *MashlingProject) InstallDependency(depPath string, version string) error {
	// Deprecated, dependency managements responsibility
	return nil
}

// Deprecated, dependency managements responsibility
func (e *MashlingProject) UninstallDependency(depPath string) error {
	// Deprecated, dependency managements responsibility
	return nil
}

func (e *MashlingProject) Build() error {
	exists := fgutil.ExecutableExists("go")
	if !exists {
		return errors.New("go not installed")
	}

	cmd := exec.Command("go", "install", "./...")
	cmd.Dir = e.GetAppDir()
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", e.GetRootDir()))
	os.Unsetenv("GOBIN")
	if e.GetDockerBuild() {
		fmt.Println("Setting GOOS to linux because this is a docker build")
		newEnv = append(newEnv, "GOOS=linux")
	}
	cmd.Env = newEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*config.FlogoAppDescriptor, error) {
	descriptor := &config.FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

//Restores dependecies using a manifest in the current working directory.
func (e *MashlingProject) RestoreDependency(manifest io.Reader) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//if manifest doesn't exists, no need to restore. return with error
	if manifest == nil {
		return ErrorNoManifest
	}

	var manifestcopy = filepath.Join(e.VendorDir, "manifest")

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
