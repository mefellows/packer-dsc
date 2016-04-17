#
# DSC Runner.
#
# Bootstraps the DSC environment, sets up configuration data
# and runs the DSC Configuration.
#
# TODO: Add parameters...?

# Set the local PowerShell Module environment path
$absoluteModulePaths = [string]::Join(";", ("c:\vagrant\modules".Split(";") | ForEach-Object { $_ | Resolve-Path }))

echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("c:\vagrant\modules".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })

$script = $(Join-Path "c:\vagrant" "manifests/Beanstalk.ps1" -Resolve)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"

# Generate the MOF file, only if a MOF path not already provided.
# Import the Manifest
. $script

cd "c:\vagrant"
$StagingPath = $(Join-Path "c:\vagrant" "staging")
$Config = $(iex (Get-Content (Join-Path "c:\vagrant" "manifests/BeanstalkConfig.psd1" -Resolve) | Out-String))
Beanstalk -OutputPath $StagingPath -MachineName "localhost" -WebAppPath "%SystemDrive%\inetpub\wwwroot" -HostName "beanstalk.dev" -ConfigurationData $Config

# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath
del $StagingPath\*.mof

# TODO: Cleanup
#rmdir c:\vagrant
