package env

type Project interface {

	// Init initializes the project settings an validates it requirements
	Init(path string) error

	// Create the project directory and its structure, optional existing vendor dir to copy
	Create(withBinDir bool, vendorDir string) error

	// Open the project directory and validate its structure
	Open() error

	// GetBinDir get the bin directory of the project
	GetBinDir() string

	// GetRootDir get the root directory of the project
	GetRootDir() string

	// GetSourceDir get the source directory of the project
	GetSourceDir() string

	// GetVendorDir get the vendor directory of the project
	GetVendorDir() string

	// GetVendorSrcDir get the vendor source directory of the project
	GetVendorSrcDir() string

	// Install a go dependency
	InstallDependency(path string, version string) error

	// Restore go dependencies
	RestoreDependency() error

	// Uninstall a go dependency
	UninstallDependency(path string) error

	// Build the project
	Build() error
}
