# Packer DSC Provisioner

A [Desired State Configuration](http://technet.microsoft.com/en-au/library/dn249912.aspx) provisioner for [Packer.io](http://packer.io), the machine image automation tool, allowing you to automate the generation of your Windows machine images in a repeatable, reliable way.

Works nicely when combined with a Vagrant development workflow, possibly leveraging the [Vagrant DSC](https://github.com/mefellows/vagrant-dsc) plugin.

Features:
* Automatically run DSC on a remote machine, from local DSC Configurations.
* Automatically configure PowerShell Package Management (specify `install_package_management` as `true`)
* Install DSC Resources from local or remote sources, including the PowerShell Gallery (see `resource_paths` and `install_modules`)
* Ability to use pre-generated MOF files, if required

[![Coverage Status](https://coveralls.io/repos/github/mefellows/packer-dsc/badge.svg?branch=HEAD)](https://coveralls.io/github/mefellows/packer-dsc?branch=HEAD)
[![wercker status](https://app.wercker.com/status/ef7336f65a3636531141a653e775d58f/s "wercker status")](https://app.wercker.com/project/bykey/ef7336f65a3636531141a653e775d58f)

### Getting Started

The plugin can be used by downloading the pre-built binary, or by building the project locally and ensuring the binary is installed in the correct location.

### On Mac OSX using Homebrew

If you are using [Homebrew](http://brew.sh) you can follow these steps to install the plugins:

```
brew install https://raw.githubusercontent.com/mefellows/packer-dsc/master/scripts/packer-provisioner-dsc.rb
```

### Using pre-built binaries

1. Install Packer
1. Download the latest [release](https://github.com/mefellows/packer-dsc/releases) for your host environment
1. Unzip the plugin binaries to [a location where Packer will detect them at run-time](https://packer.io/docs/extend/plugins.html), such as any of the following:
  - The directory where the packer binary is.
  - `~/.packer.d/plugins` on Unix systems or `%APPDATA%/packer.d/plugins` on Windows.
  - The current working directory.
1. Change to a directory where you have packer templates, and run as usual.

### Using a local build

With [Go 1.2+](http://golang.org) installed, follow these steps to use these community plugins for Windows:

1. Install packer
1. Clone this repo
1. Run `make dev`
1. Copy the plugin binaries located in `./bin` to [a location where Packer will detect them at run-time](https://packer.io/docs/extend/plugins.html), such as any of the following:
  - The directory where the packer binary is. If you've built Packer locally, then Packer and the new plugins are already in `$GOPATH/bin` together.
  - `~/.packer.d/plugins` on Unix systems or `%APPDATA%/packer.d/plugins` on Windows.
  - The current working directory.
1. Change to a directory where you have packer templates, and run as usual.

# Introduction

Type: `dsc`

DSC Packer provisioner configures DSC to run on the
machines by Packer from local modules and manifest files. Modules and manifests
can be uploaded from your local machine to the remote machine.

-&gt; **Note:** DSC will *not* be installed automatically by this
provisioner. This provisioner expects that DSC is already [installed](https://www.penflip.com/powershellorg/the-dsc-book/blob/master/dsc-overview-and-requirements.txt) on the
machine. It is common practice to use the [powershell
provisioner](/docs/provisioners/powershell.html) before the DSC provisioner to install WMF4.0+, any
required DSC Resources or connection to a DSC Pull server.

## Basic Example

The example below is fully functional and expects the configured manifest file
to exist relative to your working directory:

``` {.javascript}
{
  "type": "dsc",
  "manifest_file": "manifests/Beanstalk.ps1",
  "configuration_file": "manifests/BeanstalkConfig.psd1",
  "configuration_params": {
    "-WebAppPath": "c:\\tmp",
    "-MachineName": "localhost"
  }
}
```

## Configuration Reference

The reference of available configuration options is listed below.

Required parameters:

-   `manifest_file` (string) -  The main DSC manifest file to apply to kick off the entire thing.

Optional parameters:

-   `configuration_name` (string) -  The name of the Configuration module. Defaults to the base
    name of the `manifest_file`. e.g. `Default.ps1` would result in `Default`.

-   `mof_path` (string) -  Relative path to a folder, containing the pre-generated MOF file.

-   `configuration_file` (string) -  Relative path to the DSC Configuration Data file.
    Configuration data is used to parameterise the configuration_file.

-   `configuration_params` (object of key/value strings) - Set of Parameters to pass to the DSC
     Configuration.

-   `module_paths` (array of strings) -  Set of relative module paths.
     These paths are added to the DSC Configuration running environment to enable _local_ modules to be addressed.

-   `resource_paths` (array of strings) -  Set of DSC Resources to upload for system-wide use.
    These paths are uploaded into `${env:programfiles}\WindowsPowershell\Modules` to be used system-wide, unlike
    `module_paths` which is scoped to the current Configuration.

-   `install_modules` (array of strings) - Set of PowerShell modules to be installed
    with the `Install-Module` command. See `install_package_management` if you would
    like the DSC Provisioner to install this command for you.

-   `install_package_management` (bool) - Automatically installs the
    [Package Management](https://github.com/OneGet/oneget) package manager
    (formerly OneGet) on the server.    

-   `staging_dir` (string) - The directory where files will be uploaded.
    Packer requires write  permissions in this directory.

-   `clean_staging_dir` (bool) - If true, staging directory is removed after executing DSC.

-   `working_dir` (string) - The directory from which the command will be executed.
    Packer requires the directory to exist when running DSC.

-   `ignore_exit_codes` (boolean) - If true, Packer will never consider the
     DSC provisioning process a failure.

-   `execute_command` (string) -  The command used to execute DSC. This has
    various [configuration template
    variables](/docs/templates/configuration-templates.html) available. See
    below for more information.

## Execute Command

By default, Packer uses the following command to execute DSC:

```
{{if ne .ModulePath ""}}
$absoluteModulePaths = [string]::Join(";", ("{{.ModulePath}}".Split(";") | ForEach-Object { $_ | Resolve-Path }))
echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("{{.ModulePath}}".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })
{{end}}

$script = $("{{.ManifestFile}}" | Resolve-Path)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"

{{if eq .MofPath ""}}
# Generate the MOF file, only if a MOF path not already provided.
# Import the Manifest
. $script

cd "{{.WorkingDir}}"
$StagingPath = $(Join-Path "{{.WorkingDir}}" "staging")
{{if ne .ConfigurationFilePath ""}}
$Config = $(iex (Get-Content ("{{.WorkingDir}}" "{{.ConfigurationFilePath}}" | Resolve-Path) | Out-String))
{{end}}
{{.ConfigurationName}} -OutputPath $StagingPath {{.ConfigurationParams}}{{if ne .ConfigurationFilePath ""}} -ConfigurationData $Config{{end}}
{{else}}
$StagingPath = "{{.MofPath}}"
{{end}}

# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`
```

This command can be customized using the `execute_command` configuration. As you
can see from the default value above, the value of this configuration can
contain various template variables, defined below:

-   `WorkingDir` - The path from which DSC will be executed.
-   `ConfigurationParams` - Arguments to the DSC Configuration in k/v pairs.
-   `ConfigurationFilePath` - The path to a DSC Configuration File, if any.
-   `ConfigurationName` - The name of the DSC Configuration to run.
-   `ManifestFile` - The path on the remote machine to the manifest file for
    DSC to use.
-   `ModulePath` - The path to a directory on the remote machine containing the manifest files.
-   `MofPath` - The path to a directory containing any existing MOF file(s) to use.

## Examples

See the [Examples](examples) directory.
