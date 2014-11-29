package dsc

// The path relative to `dsc.manifests_path` pointing to the Configuration file
// dsc.configuration_file  = "MyWebsite.ps1"
//
// # The Configuration Command to run. Assumed to be the same as the `dsc.configuration_file`
// # (sans extension) if not provided.
// dsc.configuration_name = "MyWebsite"
//
// # Commandline arguments to the Configuration run
// # Set of Parameters to pass to the DSC Configuration.
// #
// # To pass in flags, simply set the value to `nil`
// dsc.configuration_params = {"machineName" => "localhost", "-EnableDebug" => nil}
//
// # Relative path to a folder containing a pre-generated MOF file.
// #
// # Path is relative to the folder containing the Vagrantfile.
// #dsc.mof_path = "mof_output"
//
// # Relative path to the folder containing the root Configuration manifest file.
// # Defaults to 'manifests'.
// #
// # Path is relative to the folder containing the Vagrantfile.
// # dsc.manifests_path = "manifests"
//
// # Set of module paths relative to the Vagrantfile dir.
// #
// # These paths are added to the DSC Configuration running
// # environment to enable local modules to be addressed.
// #
// # @return [Array] Set of relative module paths.
// #dsc.module_path = ["manifests", "modules"]
//
// # The type of synced folders to use when sharing the data
// # required for the provisioner to work properly.
// #
// # By default this will use the default synced folder type.
// # For example, you can set this to "nfs" to use NFS synced folders.
// #dsc.synced_folder_type = ""
//
// # Temporary working directory on the guest machine.
// #dsc.temp_dir = "/tmp/vagrant-dsc"
