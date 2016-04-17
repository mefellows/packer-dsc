---
description: |
DSC (local) Packer provisioner configures DSC to run on the
machines by Packer from local modules and manifest files. Modules and manifests
can be uploaded from your local machine to the remote machine.
layout: docs
page_title: 'DSC (Local) Provisioner'
...

# DSC Provisioner

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

-   `configuration_name` (string) -  The name of the Configuration module. Defaults to the base name of
    the `manifest_file`. e.g. `Default.ps1` would result in `Default`.

-   `mof_path` (string) -  Relative path to a folder, containing the pre-generated MOF file.

-   `configuration_file` (string) -  Relative path to the DSC Configuration Data file.
    Configuration data is used to parameterise the configuration_file.

-   `configuration_params` (object of key/value strings) - Set of Parameters to pass to the DSC Configuration.

-   `module_paths` (array of strings) -  Set of relative module paths.
     These paths are added to the DSC Configuration running environment to enable local modules to be addressed.

-   `staging_dir` (string) - The directory where files will be uploaded.
    Packer requires write  permissions in this directory.

-   `clean_staging_dir` (bool) - If true, staging directory is removed after executing DSC.

-   `working_dir` (string) - The directory from which the command will be executed.
    Packer requires the directory to exist when running DSC.

-   `ignore_exit_codes` (boolean) - If true, Packer will never consider the
provisioner a failure.

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
-  
