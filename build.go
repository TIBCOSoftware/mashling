package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// ImportPath is the root import path regardless of location.
	ImportPath = "github.com/TIBCOSoftware/mashling"
	// Git is the git executable.
	Git = "git"
	// Go is the go executable.
	Go = "go"
	// GoFmt is the gofmt executable.
	GoFmt = "gofmt"
	// GoBinData represents the go-bindata executable.
	GoBinData = "go-bindata"
	// GoDep represents the dep executable.
	GoDep = "dep"
)

var (
	// Platforms are the supported platforms this build process supports.
	Platforms = [...]Platform{
		{"darwin", "amd64", ""},
		{"linux", "amd64", ""},
		{"linux", "arm64", "noasm"},
		{"windows", "amd64", ""},
	}
	// Date is the date.
	Date = time.Now().Format("2006-01-02T15:04:05-0700")
	// Version is the default build version to use if no other is found.
	Version = "v0.4.0"
	// V represents verbosity for the build context.
	V = "0"
	// Env represents the environment for the build context.
	Env []string
	// TargetOS is the target platform OS for the build context.
	TargetOS string
	// TargetArch is the target platform architecture for the build context.
	TargetArch string
	// UPX is the UPX executable
	UPX string

	dependencies = make(map[string]bool)
)

// Platform represents a golang OS and ARCH build target.
type Platform struct {
	os, arch, tags string
}

func init() {
	var output []byte
	_, err := exec.LookPath(Git)
	if err != nil {
		fmt.Println("Git not found...")
	} else {
		cmd := exec.Command(Git, "describe", "--tags", "--always", "--dirty", "--match=v*")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println("Project is not a git repository...")
		}
	}
	if len(output) == 0 {
		output, err = ioutil.ReadFile(".version")
		if err != nil {
			fmt.Println(".version file not found, using default version:", Version)
		}
	}
	if len(output) > 0 {
		Version = strings.TrimSpace(string(output))
	}

	Env = os.Environ()
	Env = append(Env, "CGO_ENABLED=0")

	v := os.Getenv("V")
	if v != "" {
		V = v
	}
}

func GetCurrentPlatform() (Platform, error) {
	os := TargetOS
	arch := TargetArch
	if os == "" {
		os = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	for _, platform := range Platforms {
		if platform.os == os && platform.arch == arch {
			return platform, nil
		}
	}

	return Platform{}, errors.New("platform not found")
}

// Step is a function that represents a step in the build process.
type Step func() error

// Resolve runs a step function.
func Resolve(step Step) {
	name := fmt.Sprintf("%#v", Step(step))
	if !dependencies[name] {
		err := step()
		if err != nil {
			panic(err)
		}
		dependencies[name] = true
	}
}

// Print outputs a message and mimics the old Makefile output.
func Print(message string) {
	fmt.Println("\033[34;1m▶\033[0m", message)
}

// Dir returns the path of a file.
func Dir(file string) string {
	parts, path := strings.Split(ImportPath, string(os.PathSeparator)), ""
	for _, part := range parts[:len(parts)-1] {
		path += part + "/"
	}
	return path
}

// List lists all go files within the build context and applies a function.
func List(path string, processFiles func(files []string) error) error {
	// path = filepath.Join(Primary, "src", ImportPath, path)
	var process func(path string) error
	process = func(path string) error {
		dir, err := os.Open(path)
		if err != nil {
			return err
		}
		var files []string
		entries, err := dir.Readdir(256)
		for err == nil {
			for _, entry := range entries {
				if entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".") || entry.Mode()&os.ModeSymlink != 0 {
					continue
				} else if entry.IsDir() {
					er := process(filepath.Join(path, entry.Name()))
					if er != nil {
						return er
					}
				} else if strings.HasSuffix(entry.Name(), ".go") {
					files = append(files, filepath.Join(path, entry.Name()))
				}
			}
			entries, err = dir.Readdir(256)
		}
		if len(files) > 0 {
			return processFiles(files)
		}
		return nil
	}
	return process(path)
}

// Exec executes a command within the build context.
func Exec(name string, arg ...string) error {
	if V == "1" {
		fmt.Println(name + " " + strings.Join(arg, " "))
	}
	cmd := exec.Command(name, arg...)
	cmd.Env = Env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}

	return nil
}

// Ldflags returns default ldflags for the build context.
func Ldflags() string {
	return fmt.Sprintf("-X %s/internal/app/version.Version=%s -X %s/internal/app/version.BuildDate=%s", ImportPath, Version, ImportPath, Date)
}

// ReleaseLdflags returns flags required for a release build.
func ReleaseLdflags() string {
	return fmt.Sprintf("-s -w %s", Ldflags())
}

func main() {
	buildSet := flag.NewFlagSet("", flag.ExitOnError)
	buildSet.StringVar(&TargetOS, "os", "", "target OS")
	buildSet.StringVar(&TargetArch, "arch", "", "target ARCH")
	if len(os.Args) > 2 {
		buildSet.Parse(os.Args[2:])
	}
	args := os.Args

	if TargetOS != "" {
		if TargetArch == "" {
			TargetArch = "amd64"
		}
		found := false
		for _, platform := range Platforms {
			if platform.os == TargetOS && platform.arch == TargetArch {
				found = true
				break
			}
		}
		if !found {
			Print(fmt.Sprintf("Invalid OS and Arch combination: %s - %s", TargetOS, TargetArch))
			return
		}
	}
	if len(args) > 1 {
		switch args[1] {
		case "build":
			Resolve(build)
		case "buildgateway":
			Resolve(buildgateway)
		case "buildcli":
			Resolve(buildcli)
		case "all":
			Resolve(all)
		case "allgateway":
			Resolve(allgateway)
		case "allgatewayprep":
			Resolve(allgatewayprep)
		case "allcli":
			Resolve(allcli)
		case "allcliprep":
			Resolve(allcliprep)
		case "releaseall":
			Resolve(releaseall)
		case "releasegateway":
			Resolve(releasegateway)
		case "releasecli":
			Resolve(releasecli)
		case "docker":
			Resolve(docker)
		case "setup":
			Resolve(setup)
		case "gofmt":
			Resolve(gofmt)
		case "vet":
			Resolve(vet)
		case "generate":
			Resolve(generate)
		case "cligenerate":
			Resolve(cligenerate)
		case "gatewayassets":
			Resolve(gatewayassets)
		case "cliassets":
			Resolve(cliassets)
		case "clean":
			Resolve(clean)
		case "version":
			Resolve(version)
		default:
			Print(fmt.Sprintf("Unrecognized command: %s", args[1]))
		}
		return
	}
	Resolve(build)
}

func build() error {
	Resolve(buildgateway)
	Resolve(buildcli)
	return nil
}

func buildgateway() error {
	Print("building gateway executable...")

	arg := []string{"install"}
	platform, err := GetCurrentPlatform()
	if err == nil && platform.tags != "" {
		arg = append(arg, "-tags", platform.tags)
	}
	arg = append(arg, "-ldflags", Ldflags(), fmt.Sprintf("%s/cmd/mashling-gateway", ImportPath))
	err = Exec(Go, arg...)
	if err != nil {
		return err
	}

	return nil
}

func buildcli() error {
	Print("building CLI executable...")

	arg := []string{"install"}
	platform, err := GetCurrentPlatform()
	if err == nil && platform.tags != "" {
		arg = append(arg, "-tags", platform.tags)
	}
	arg = append(arg, "-ldflags", Ldflags(), fmt.Sprintf("%s/cmd/mashling-cli", ImportPath))
	err = Exec(Go, arg...)
	if err != nil {
		return err
	}

	return nil
}

func all() error {
	Resolve(allgateway)
	Resolve(allcli)
	return nil
}

func allgateway() error {
	Resolve(allgatewayprep)
	Resolve(vet)
	Resolve(buildgateway)

	return nil
}

func allgatewayprep() error {
	Resolve(gatewayassets)
	Resolve(generate)
	Resolve(gofmt)
	return nil
}

func allcli() error {
	Resolve(allcliprep)
	Resolve(vet)
	Resolve(buildcli)

	return nil
}

func allcliprep() error {
	Resolve(cligenerate)
	Resolve(cliassets)
	Resolve(gofmt)

	return nil
}

func releaseall() error {
	Resolve(platforms)
	Resolve(upx)

	Print("building release executables")

	for _, platform := range Platforms {
		gateway, err := releaseGatewayWithTarget(platform.os, platform.arch, platform.tags)
		if err != nil {
			return err
		}
		cli, cErr := releaseCLIWithTarget(platform.os, platform.arch, platform.tags)
		if cErr != nil {
			return cErr
		}

		if UPX != "" {
			err = Exec(UPX, gateway, cli)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func releasegateway() error {
	Resolve(platforms)
	Resolve(upx)

	Print("building gateway release executable")
	platform, err := GetCurrentPlatform()
	if err != nil {
		return err
	}
	gateway, err := releaseGatewayWithTarget(platform.os, platform.arch, platform.tags)
	if err != nil {
		return err
	}
	if UPX != "" {
		err = Exec(UPX, gateway)
		if err != nil {
			return err
		}
	}

	return nil
}

func releasecli() error {
	Resolve(platforms)
	Resolve(upx)

	Print("building CLI release executable")
	platform, err := GetCurrentPlatform()
	if err != nil {
		return err
	}
	cli, err := releaseCLIWithTarget(platform.os, platform.arch, platform.tags)
	if err != nil {
		return err
	}
	if UPX != "" {
		err = Exec(UPX, cli)
		if err != nil {
			return err
		}
	}

	return nil
}

func releaseGatewayWithTarget(os string, arch string, tags string) (string, error) {
	var extension string
	if os == "windows" {
		extension = ".exe"
	}
	gateway := fmt.Sprintf("release/mashling-gateway-%s-%s%s", os, arch, extension)
	cmd := exec.Command(Go, "build", "-tags", "release "+tags,
		"-ldflags", ReleaseLdflags(),
		"-o", gateway,
		fmt.Sprintf("%s/cmd/mashling-gateway", ImportPath))
	cmd.Env = append(Env,
		fmt.Sprintf("GOOS=%v", os),
		fmt.Sprintf("GOARCH=%s", arch))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return gateway, fmt.Errorf("%v: %s", err, string(output))
	}
	return gateway, nil
}

func releaseCLIWithTarget(os string, arch string, tags string) (string, error) {
	var extension string
	if os == "windows" {
		extension = ".exe"
	}
	cli := fmt.Sprintf("release/mashling-cli-%s-%s%s", os, arch, extension)
	cmd := exec.Command(Go, "build", "-tags", "release "+tags,
		"-ldflags", ReleaseLdflags(),
		"-o", cli,
		fmt.Sprintf("%s/cmd/mashling-cli", ImportPath))
	cmd.Env = append(Env,
		fmt.Sprintf("GOOS=%v", os),
		fmt.Sprintf("GOARCH=%s", arch))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return cli, fmt.Errorf("%v: %s", err, string(output))
	}
	return cli, nil
}

func platforms() error {

	return nil
}

func gofmt() error {
	Print("running gofmt…")

	err := List(".", func(files []string) error {
		args := append([]string{"-l", "-w"}, files...)
		err := Exec(GoFmt, args...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func vet() error {
	Print("running go vet…")

	err := List(".", func(files []string) error {
		args := append([]string{"vet"}, files...)
		err := Exec(Go, args...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func generate() error {
	Print("running go generate…")

	err := Exec(Go, "generate", "./...")
	if err != nil {
		return err
	}

	return nil
}

func cligenerate() error {
	Print("running CLI go generate…")

	err := Exec(Go, "generate", "./internal/app/cli/...")
	if err != nil {
		return err
	}

	return nil
}

func gatewayassets() error {
	Resolve(gobindata)

	Print("running asset generation…")

	sets := make(map[string][]string)
	process := func(fileType, path, message string) {
		source := strings.Replace(path, ".json", ".go", -1)
		_, err := os.Stat(source)
		if err != nil {
			return
		}
		fmt.Println(message, source)
		set := sets[fileType]
		set = append(set, path)
		sets[fileType] = set
	}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		switch name := info.Name(); name {
		case "activity.json":
			process(name, path, "Activity Found:")
		case "trigger.json":
			process(name, path, "Trigger Found:")
		case "action.json":
			process(name, path, "Asset Found:")
		}
		return nil
	})
	if err != nil {
		return err
	}

	var assets []string
	err = filepath.Walk("internal/app/assets/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasSuffix(path, ".go") {
			return nil
		}
		assets = append(assets, path)
		return nil
	})
	if err != nil {
		return err
	}

	if activities := sets["activity.json"]; len(activities) > 0 {
		err = Exec(GoBinData, append([]string{"-pkg", "activities", "-o",
			"internal/app/gateway/flogo/registry/activities/activities.go"},
			activities...)...)
		if err != nil {
			return err
		}
	}

	if triggers := sets["trigger.json"]; len(triggers) > 0 {
		err = Exec(GoBinData, append([]string{"-pkg", "triggers", "-o",
			"internal/app/gateway/flogo/registry/triggers/triggers.go"},
			triggers...)...)
		if err != nil {
			return err
		}
	}

	if actions := sets["action.json"]; len(actions) > 0 {
		err = Exec(GoBinData, append([]string{"-pkg", "actions", "-o",
			"internal/app/gateway/flogo/registry/actions/actions.go"},
			actions...)...)
		if err != nil {
			return err
		}
	}

	err = Exec(GoBinData, append([]string{"-prefix",
		"internal/app/assets/", "-pkg", "assets", "-o",
		"internal/app/assets/assets.go"},
		assets...)...)
	if err != nil {
		return err
	}

	return nil
}

func cliassets() error {
	Resolve(gobindata)

	Print("running CLI asset generation…")

	var assets []string
	err := filepath.Walk("internal/app/cli/assets/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasSuffix(path, ".go") {
			return nil
		}
		assets = append(assets, path)
		return nil
	})
	if err != nil {
		return err
	}

	err = Exec(GoBinData, append([]string{"-prefix",
		"internal/app/cli/assets/", "-pkg", "assets", "-o",
		"internal/app/cli/assets/assets.go"},
		assets...)...)
	if err != nil {
		return err
	}

	return nil
}

func setup() error {
	Resolve(clean)
	Resolve(godep)
	Resolve(gobindata)
	return nil
}

func clean() error {
	Print("cleaning...")
	err := os.RemoveAll("./release")
	if err != nil {
		return err
	}
	return nil
}

func version() error {
	Print(Version)
	return nil
}

func docker() error {
	Resolve(platforms)
	TargetOS = "linux"
	TargetArch = "amd64"
	Resolve(releasegateway)

	Print("building a docker image containing the mashling-gateway binary")

	docker, err := exec.LookPath("docker")
	if err != nil {
		return errors.New("docker not found")
	}

	err = Exec(docker, "build", "-f", "dockerfiles/run/Dockerfile", ".", "-t", "mashling-gateway")
	if err != nil {
		return err
	}

	return nil
}

func gobindata() error {
	_, err := exec.LookPath(GoBinData)
	if err == nil {
		return nil
	}

	Print("building go-bindata…")

	err = Exec(Go, "get", "github.com/mashling-support/go-bindata/...")
	if err != nil {
		return err
	}

	return nil
}

func godep() error {
	_, err := exec.LookPath(GoDep)
	if err == nil {
		return nil
	}

	Print("building dep…")

	err = Exec(Go, "get", "github.com/golang/dep/cmd/dep")
	if err != nil {
		return err
	}

	return nil
}

func upx() error {
	var err error
	UPX, err = exec.LookPath("upx")
	if err != nil {
		UPX = ""
		fmt.Println("UPX not found, skipping compression (please visit https://upx.github.io to install)...")
	}
	return nil
}
