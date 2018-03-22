package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/callum-ramage/jsonconfig"
	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/dep"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"go/build"
)

// dockerfile is the template for a dockerfile needed to build a docker image
const dockerfile = `# Dockerfile for {{.name}}
# VERSION {{.version}}
FROM alpine
RUN apk update && apk add ca-certificates
ADD linux_amd64 .
EXPOSE {{.port}}
CMD ./linux_amd64/{{.name}}`

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env env.Project) error
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson, appDir, appName, vendorDir, constraints string) error {
	return doCreate(env, appJson, appDir, appName, vendorDir, constraints)
}

// doCreate performs the app creation
func doCreate(enviro env.Project, appJson, rootDir, appName, vendorDir, constraints string) error {
	fmt.Printf("Creating initial project structure, this might take a few seconds ... \n")
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}
	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	} else {
		appName = descriptor.Name
		rootDir = filepath.Join(rootDir, appName)
	}

	err = enviro.Init(rootDir)
	if err != nil {
		return err
	}

	err = enviro.Create(false, "")
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(filepath.Join(rootDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}
	// create initial structure
	appDir := filepath.Join(enviro.GetSourceDir(), descriptor.Name)
	os.MkdirAll(appDir, os.ModePerm)

	// Validate structure
	err = enviro.Open()
	if err != nil {
		return err
	}

	// Create the dep manager
	depManager := &dep.DepManager{Env: enviro}

	// Initialize the dep manager
	err = depManager.Init()
	if err != nil {
		return err
	}

	// Create initial files
	deps, err := config.ExtractAllDependencies(appJson)
	if err != nil {
		return err
	}

	createMainGoFile(appDir, "")
	createImportsGoFile(appDir, deps)

	// Add constraints
	if len(constraints) > 0 {
		newConstraints := []string{"-add"}
		newConstraints = append(newConstraints, strings.Split(constraints, ",")...)
		err = depManager.Ensure(newConstraints...)
		if err != nil {
			return err
		}
	}

	ensureArgs := []string{}

	if len(vendorDir) > 0 {
		// Copy vendor directory
		fgutil.CopyDir(vendorDir, enviro.GetVendorDir())
		// Do not touch vendor folder when ensuring
		ensureArgs = append(ensureArgs, "-no-vendor")
	}

	// Sync up
	err = depManager.Ensure(ensureArgs...)
	if err != nil {
		return err
	}

	return nil
}

type PrepareOptions struct {
	PreProcessor    BuildPreProcessor
	OptimizeImports bool
	EmbedConfig     bool
	Shim            string
}

// PrepareApp do all pre-build setup and pre-processing
func PrepareApp(env env.Project, options *PrepareOptions) error {
	return doPrepare(env, options)
}

// doPrepare performs all the prepare functionality
func doPrepare(env env.Project, options *PrepareOptions) (err error) {
	// Create the dep manager
	depManager := dep.DepManager{Env: env}
	if !depManager.IsInitialized() {
		// This is an old app
		err = MigrateOldApp(env, depManager)
		if err != nil {
			return err
		}
	}

	if options == nil {
		options = &PrepareOptions{}
	}

	// Call external preprocessor
	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(env)
		if err != nil {
			return err
		}
	}

	//generate metadata
	err = generateGoMetadata(env)
	if err != nil {
		return err
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(filepath.Join(env.GetRootDir(), "flogo.json"))

	if err != nil {
		return err
	}
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	removeEmbeddedAppGoFile(env.GetAppDir())
	removeShimGoFiles(env.GetAppDir())

	if options.Shim != "" {

		removeMainGoFile(env.GetAppDir()) //todo maybe rename if it exists
		createShimSupportGoFile(env.GetAppDir(), appJson, options.EmbedConfig)

		fmt.Println("Shim:", options.Shim)

		for _, value := range descriptor.Triggers {

			fmt.Println("Id:", value.Id)
			if value.Id == options.Shim {
				triggerPath := filepath.Join(env.GetVendorSrcDir(), value.Ref, "trigger.json")

				mdJson, err := fgutil.LoadLocalFile(triggerPath)
				if err != nil {
					return err
				}
				metadata, err := ParseTriggerMetadata(mdJson)
				if err != nil {
					return err
				}

				fmt.Println("Shim Metadata:", metadata.Shim)

				if metadata.Shim != "" {

					//todo blow up if shim file not found
					shimFilePath := filepath.Join(env.GetVendorSrcDir(), value.Ref, dirShim, fileShimGo)
					fmt.Println("Shim File:", shimFilePath)
					fgutil.CopyFile(shimFilePath, filepath.Join(env.GetAppDir(), fileShimGo))
					
					// ensure deps after the shim.go has been copied to main.go...
					depManager.Ensure()
					
					// This is a bit of a workaround, will resolve with a better solution in the future
					// generate metadata again... ensure will remove it
					err = generateGoMetadata(env)
					if err != nil {
						return err
					}					

					if metadata.Shim == "plugin" {
						//look for Makefile and execute it
						makeFilePath := filepath.Join(env.GetVendorSrcDir(), value.Ref, dirShim, makeFile)
						fmt.Println("Make File:", makeFilePath)
						fgutil.CopyFile(makeFilePath, filepath.Join(env.GetAppDir(), makeFile))

						// Execute make
						cmd := exec.Command("make", "-C", env.GetAppDir())
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("GOPATH=%s", env.GetRootDir()),
						)

						err = cmd.Run()
						if err != nil {
							return err
						}
					}
				}

				break
			}
		}

	} else if options.EmbedConfig {
		createEmbeddedAppGoFile(env.GetAppDir(), appJson)
	}
	return
}

func MigrateOldApp(env env.Project, depManager dep.DepManager) error {
	// This is an old app

	// Move old vendor folder to /src/<my_app>/vendor/
	oldVendorDir := filepath.Join(env.GetRootDir(), "vendor")
	_, err := os.Stat(oldVendorDir)
	if err == nil {
		// Vendor found, move it
		err = CopyDir(oldVendorDir, env.GetVendorDir())
		if err != nil {
			return err
		}
		defer os.RemoveAll(oldVendorDir)
	}

	fmt.Println("Initializing dependency management files ....")
	err = depManager.Init()
	if err != nil {
		return err
	}

	return nil
}

type BuildOptions struct {
	*PrepareOptions

	NoGeneration   bool
	GenerationOnly bool
	SkipPrepare    bool
	BuildDocker string
}

// BuildApp build the flogo application
func BuildApp(env env.Project, options *BuildOptions) error {
	return doBuild(env, options)
}

// doBuildApp performs the build functionality
func doBuild(env env.Project, options *BuildOptions) (err error) {
	if options == nil {
		options = &BuildOptions{}
	}

	if len(options.BuildDocker) > 0 {
    		env.SetDockerBuild()
    }

	if options.GenerationOnly {
		// Only perform prepare
		return PrepareApp(env, options.PrepareOptions)
	}

	if !options.SkipPrepare && !options.NoGeneration {
		err = PrepareApp(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	if !options.EmbedConfig {
		fgutil.CopyFile(filepath.Join(env.GetRootDir(), config.FileDescriptor), filepath.Join(env.GetBinDir(), config.FileDescriptor))
		if err != nil {
			return err
		}
	} else {
		os.Remove(filepath.Join(env.GetBinDir(), config.FileDescriptor))
	}

	  // To create a dockerfile this component executes four steps
    // 1. Check if flogo.json exists in bin folder (built without -e)
    // 2. Read flogo.json from ./flogo.json
    // 3. Output the dockerfile in ./bin/dockerfile
    // 4. Execute docker build
    if len(options.BuildDocker) > 0 {
        fmt.Println("docker:", options.BuildDocker)
        config, err := jsonconfig.LoadAbstract("./flogo.json", "")
        if err != nil {
            return err
        }
        data := make(map[string]interface{})
        found := false

        for _, value := range config["triggers"].Arr {
            if value.Obj["id"].Str == options.BuildDocker {
                found = true
                data["name"] = config["name"].Str
                data["version"] = config["version"].Str
                data["port"] = value.Obj["settings.port"].Str
            }
        }

        if options.BuildDocker == "no-trigger" {
            found = true
            data["name"] = config["name"].Str
            data["version"] = config["version"].Str
            data["port"] = ""
        }

        if found {
            t := template.Must(template.New("email").Parse(dockerfile))
            buf := &bytes.Buffer{}
            if err := t.Execute(buf, data); err != nil {
                return err
            }
            s := buf.String()

            if data["port"] == "" {
                s = strings.Replace(s, "EXPOSE \n", "", -1)
            }

            file, err := os.Create("./bin/dockerfile")
            if err != nil {
                return err
            }
            defer file.Close()

            file.WriteString(s)
            file.Sync()

            cmd := exec.Command("docker", "build", ".", "-t", strings.ToLower(config["name"].Str)+":"+config["version"].Str)
            cmd.Dir = "./bin"
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            err = cmd.Run()
            if err != nil {
                return err
            }
        } else {
            fmt.Println("Your app doesn't contain the trigger you specified so we can't create a dockerfile for it")
        }
    }
	return
}

// InstallPalette install a palette
func InstallPalette(env env.Project, path string) error {

	var file []byte

	file, _ = ioutil.ReadFile(path)

	var paletteDescriptor *config.FlogoPaletteDescriptor
	err := json.Unmarshal(file, &paletteDescriptor)

	var deps []config.Dependency

	if err != nil {
		err = json.Unmarshal(file, &deps)
	} else {
		deps = paletteDescriptor.Extensions
	}

	if err != nil {
		return err
		//fmt.Fprint(os.Stderr, "Error: Unable to parse palette descriptor, file may be corrupted.\n\n")
		//os.Exit(2)
	}

	for _, dep := range deps {
		err = InstallDependency(env, dep.Ref, "")
		if err != nil {
			return err
		}
	}

	//fmt.Fprintf(os.Stdout, "Adding Palette: %s [%s]\n\n", paletteDescriptor.Name, paletteDescriptor.Description)

	return nil
}

// InstallDependency install a dependency
func InstallDependency(environ env.Project, path string, version string) error {
	// Create the dep manager
	depManager := dep.DepManager{Env: environ}
	if !depManager.IsInitialized() {
		// This is an old app
		err := MigrateOldApp(environ, depManager)
		if err != nil {
			return err
		}
	}
	err := depManager.InstallDependency(path, version)
	if err != nil {
		return err
	}
	return nil
}

// UninstallDependency uninstall a dependency
func UninstallDependency(environ env.Project, path string) error {
	// Create the dep manager
	depManager := dep.DepManager{Env: environ}
	if !depManager.IsInitialized() {
		// This is an old app
		err := MigrateOldApp(environ, depManager)
		if err != nil {
			return err
		}
	}
	err := depManager.UninstallDependency(path)
	if err != nil {
		return err
	}
	return nil
}

func ListDependencies(env env.Project, cType config.ContribType) ([]*config.Dependency, error) {
	// Get build context
	bc := build.Default
	currentGoPath := bc.GOPATH
	bc.GOPATH = env.GetRootDir()
	defer func() { bc.GOPATH = currentGoPath }()
	pkgs, err := bc.ImportDir(env.GetAppDir(), build.IgnoreVendor)
	if err != nil {
		return nil, err
	}
	var deps []*config.Dependency
	// Get all imports
	for _, imp := range pkgs.Imports {
		pkg, err := bc.Import(imp, env.GetAppDir(), build.FindOnly)
		if err != nil {
			// Ignore package
			continue
		}
		if cType == 0 || cType == config.ACTION {
			filePath := filepath.Join(pkg.Dir, "action.json")
			// Check if it is an action
			info, err := os.Stat(filePath)
			if err == nil {
				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:action" {
					deps = append(deps, &config.Dependency{ContribType: config.ACTION, Ref: imp})
				}
			}
		}
		if cType == 0 || cType == config.TRIGGER {
			filePath := filepath.Join(pkg.Dir, "trigger.json")
			// Check if it is a trigger
			info, err := os.Stat(filePath)
			if err == nil {
				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:trigger" {
					deps = append(deps, &config.Dependency{ContribType: config.TRIGGER, Ref: imp})
				}
			}
		}
		if cType == 0 || cType == config.ACTIVITY {
			filePath := filepath.Join(pkg.Dir, "activity.json")
			// Check if it is an activity
			info, err := os.Stat(filePath)
			if err == nil {
				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:activity" {
					deps = append(deps, &config.Dependency{ContribType: config.ACTIVITY, Ref: imp})
				}
			}
		}
		if cType == 0 || cType == config.FLOW_MODEL {
			filePath := filepath.Join(pkg.Dir, "flow-model.json")
			// Check if it is a flow model
			info, err := os.Stat(filePath)
			if err == nil {
				desc, err := readDescriptor(filePath, info)
				if err == nil && desc.Type == "flogo:flow-model" {
					deps = append(deps, &config.Dependency{ContribType: config.FLOW_MODEL, Ref: imp})
				}
			}
		}
	}
	return deps, nil
}

// Ensure is a wrapper for dep ensure command
func Ensure(depManager *dep.DepManager, args ...string) error{
	return depManager.Ensure(args...)
}

func readDescriptor(path string, info os.FileInfo) (*config.Descriptor, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return ParseDescriptor(string(raw))
}

func generateGoMetadata(env env.Project) error {
	//todo optimize metadata recreation to minimize compile times
	dependencies, err := ListDependencies(env, 0)

	if err != nil {
		return err
	}

	for _, dependency := range dependencies {
		createMetadata(env, dependency)
	}

	return nil
}

func createMetadata(env env.Project, dependency *config.Dependency) error {

	vendorSrc := env.GetVendorSrcDir()
	mdFilePath := filepath.Join(vendorSrc, dependency.Ref)
	mdGoFilePath := filepath.Join(vendorSrc, dependency.Ref)
	pkg := filepath.Base(mdFilePath)

	tplMetadata := tplMetadataGoFile

	switch dependency.ContribType {
	case config.ACTION:
		mdFilePath = filepath.Join(mdFilePath, "action.json")
		mdGoFilePath = filepath.Join(mdGoFilePath, "action_metadata.go")
	case config.TRIGGER:
		mdFilePath = filepath.Join(mdFilePath, "trigger.json")
		mdGoFilePath = filepath.Join(mdGoFilePath, "trigger_metadata.go")
		tplMetadata = tplTriggerMetadataGoFile
	case config.ACTIVITY:
		mdFilePath = filepath.Join(mdFilePath, "activity.json")
		mdGoFilePath = filepath.Join(mdGoFilePath, "activity_metadata.go")
		tplMetadata = tplActivityMetadataGoFile
	default:
		return nil
	}

	raw, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: string(raw),
	}

	f, _ := os.Create(mdGoFilePath)
	fgutil.RenderTemplate(f, tplMetadata, info)
	f.Close()

	return nil
}

var tplMetadataGoFile = `package {{.Package}}

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

func getJsonMetadata() string {
	return jsonMetadata
}
`

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
`

// ParseDescriptor parse a descriptor
func ParseDescriptor(descJson string) (*config.Descriptor, error) {
	descriptor := &config.Descriptor{}

	err := json.Unmarshal([]byte(descJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
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

// ParseTriggerMetadata parse the trigger metadata
func ParseTriggerMetadata(metadataJson string) (*config.TriggerMetadata, error) {
	metadata := &config.TriggerMetadata{}

	err := json.Unmarshal([]byte(metadataJson), metadata)

	if err != nil {
		return nil, err
	}

	return metadata, nil
}
