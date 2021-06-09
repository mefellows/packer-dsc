package dsc

import (
	common "github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ctx                 interpolate.Context

	// The command used to execute DSC.
	ExecuteCommand string `mapstructure:"execute_command"`

	// The file containing the command to execute DSC
	ExecuteCommandFilePath string `mapstructure:"execute_command_file"`

	// Set of Parameters to pass to the DSC Configuration.
	ConfigurationParams map[string]string `mapstructure:"configuration_params"`

	// Relative path to a folder, containing the pre-generated MOF file.
	//
	// Path is relative to the folder containing the Packer json.
	MofPath string `mapstructure:"mof_path"`

	// Relative path to the DSC Configuration Data file.
	//
	// Configuration data is used to parameterise the configuration_file.
	//
	// Path is relative to the folder containing the Packer json.
	ConfigurationFilePath string `mapstructure:"configuration_file"`

	// Relative path to the folder containing the root Configuration manifest file.
	// Defaults to 'manifests'.
	//
	// Path is relative to the folder containing the Packer json.
	ManifestDir string `mapstructure:"manifest_dir"`

	// The main DSC manifest file to apply to kick off the entire thing.
	//
	// Path is relative to the folder containing the Packer json.
	ManifestFile string `mapstructure:"manifest_file"`

	// The name of the Configuration module
	//
	// Defaults to the basename of the "configuration_file"
	// e.g. "Foo.ps1" becomes "Foo"
	ConfigurationName string `mapstructure:"configuration_name"`

	// Set of module paths relative to the Packer json dir.
	//
	// These paths are added to the DSC Configuration running
	// environment to enable local modules to be addressed.
	ModulePaths []string `mapstructure:"module_paths"`

	// Set of DSC resources to upload for system-wide use.
	//
	// These paths are uploaded into %SystemDrive%\WindowsPowershell\Modules
	// to be used system-wide.
	ResourcePaths []string `mapstructure:"resource_paths"`

	// Install the latest Windows PackageManagement software?
	InstallPackageManagement bool `mapstructure:"install_package_management"`

	// Modules to install, using the latest PackageManagement tooling
	// e.g. { "xWebAdministration": "1.0.0.0" }
	//
	// See InstallPackageManagement if
	InstallModules map[string]string `mapstructure:"install_modules"`

	// The directory where files will be uploaded. Packer requires write
	// permissions in this directory.
	StagingDir string `mapstructure:"staging_dir"`

	// If true, staging directory is removed after executing dsc.
	CleanStagingDir bool `mapstructure:"clean_staging_dir"`

	// The directory from which the command will be executed.
	// Packer requires the directory to exist when running dsc.
	WorkingDir string `mapstructure:"working_dir"`

	// If true, packer will ignore all exit-codes from a dsc run
	IgnoreExitCodes bool `mapstructure:"ignore_exit_codes"`

	// Specify remote DSC resources to be installed prior to the DSC execution
	// InstallResources map[string]string  `mapstructure:"install_resources"`
}
