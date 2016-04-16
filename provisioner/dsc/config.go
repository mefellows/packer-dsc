package dsc

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ctx                 interpolate.Context

	// The command used to execute Puppet.
	ExecuteCommand string `mapstructure:"execute_command"`

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
	ConfigurationFilePath string `mapstructure:"configuration_file_path"`

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

	// The type of synced folders to use when sharing the data
	// required for the provisioner to work properly.
	//
	// By default this will use the default synced folder type.
	// For example, you can set this to "nfs" to use NFS synced folders.
	SyncedFolderType string `mapstructure:"synced_folder_type"`

	// Temporary working directory on the guest machine.
	TempDir string `mapstructure:"temp_dir"`

	// The directory where files will be uploaded. Packer requires write
	// permissions in this directory.
	StagingDir string `mapstructure:"staging_directory"`

	// If true, staging directory is removed after executing dsc.
	CleanStagingDir bool `mapstructure:"clean_staging_directory"`

	// The directory from which the command will be executed.
	// Packer requires the directory to exist when running dsc.
	WorkingDir string `mapstructure:"working_directory"`

	// If true, packer will ignore all exit-codes from a dsc run
	IgnoreExitCodes bool `mapstructure:"ignore_exit_codes"`

	// Specify remote DSC resources to be installed prior to the DSC execution
	// InstallResources map[string]string  `mapstructure:"install_resources"`
}
