package dsc

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func testConfig() map[string]interface{} {
	filename := "/tmp/packer-dsc-pull-manifest"
	if _, err := os.Stat(filename); err != nil {
		tf, err := ioutil.TempFile("", "packer")
		os.Rename(tf.Name(), "/tmp/packer-dsc-pull-manifest")

		if err != nil {
			panic(err)
		}
	}

	return map[string]interface{}{
		"manifest_file":      filename,
		"manifest_dir":       ".",
		"configuration_file": "./provisioner_test.go",
		"configuration_params": map[string]string{
			"-Foo": "bar",
		},
		"resource_paths": []string{
			".",
		},
		"install_package_management": true,
		"install_modules": map[string]string{
			"SomeModule1": "1.0.0",
			"SomeModule2": "2.0.0",
		},
	}
}

func TestProvisioner_Impl(t *testing.T) {
	var raw interface{}
	raw = &Provisioner{}
	if _, ok := raw.(packer.Provisioner); !ok {
		t.Fatalf("must be a Provisioner")
	}
}

func TestProvisionerPrepare_configurationFileDataPath(t *testing.T) {
	config := testConfig()

	delete(config, "configuration_file")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test one that doesn't exist
	config["configuration_file"] = "does/not.exist"
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// And one that is a dir
	config["configuration_file"] = "."
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// Test with a good one
	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("error tempfile: %s", err)
	}
	defer os.Remove(tf.Name())

	config["configuration_file"] = tf.Name()
	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_manifestFile(t *testing.T) {
	config := testConfig()

	delete(config, "manifest_file")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test one that doesn't exist
	config["manifest_file"] = "this/file/doesnt.exist"
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// Test with a good one
	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("error tempfile: %s", err)
	}
	defer os.Remove(tf.Name())

	config["manifest_file"] = tf.Name()
	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProvisionerPrepare_manifestDir(t *testing.T) {
	config := testConfig()

	delete(config, "manifestdir")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// File is not a dir!
	delete(config, "manifestdir")
	config["manifest_dir"] = "./provisioner_test.go"
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// Dir not exist
	delete(config, "manifestdir")
	config["manifest_dir"] = "i/do/not/exist"
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// Test with a good one
	td, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	defer os.RemoveAll(td)

	config["manifest_dir"] = td
	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_modulePaths(t *testing.T) {
	config := testConfig()

	delete(config, "module_paths")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// File is not a dir
	config["module_paths"] = []string{"provisioner_test.go"}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test with bad paths
	config["module_paths"] = []string{"i-should-not-exist"}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test with a good one
	td, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	defer os.RemoveAll(td)

	config["module_paths"] = []string{td}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvisionerPrepare_resourcePaths(t *testing.T) {
	config := testConfig()

	delete(config, "resource_paths")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// File is not a dir
	config["resource_paths"] = []string{"provisioner_test.go"}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test with bad paths
	config["resource_paths"] = []string{"i-should-not-exist"}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test with a good one
	td, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	defer os.RemoveAll(td)

	config["resource_paths"] = []string{td}
	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_configurationParams(t *testing.T) {
	config := testConfig()

	delete(config, "configuration_params")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test with malformed fact
	config["configuration_params"] = "fact=stringified"
	p = new(Provisioner)
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should be an error")
	}

	// Test with a good one
	td, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	defer os.RemoveAll(td)

	facts := make(map[string]string)
	facts["fact_name"] = "fact_value"
	config["configuration_params"] = facts

	p = new(Provisioner)
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Make sure the default facts are present
	delete(config, "configuration_params")
	p = new(Provisioner)
	err = p.Prepare(config)
	if p.config.ConfigurationParams == nil {
		t.Fatalf("err: Default facts are not set in the Puppet provisioner!")
	}
}

func TestProvisionerPrepare_installPackage(t *testing.T) {
	config := testConfig()
	config["install_modules"] = map[string]string{
		"SomeModuleName":  "1.0.0",
		"SomeModuleName2": "2.0.0",
	}
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisioner_installPackage(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	config["install_modules"] = map[string]string{
		"SomeModuleName":  "1.0.0",
		"SomeModuleName2": "2.0.0",
	}
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = p.installPackage(ctx, ui, comm, "SomeModuleName", "1.0.0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expectedCommand := `powershell "& { Install-Module -Name SomeModuleName -RequiredVersion 1.0.0 -Force; exit $LastExitCode}"`
	if comm.StartCmd.Command != expectedCommand {
		t.Fatalf("Expected command '%s' but got '%s'", expectedCommand, comm.StartCmd.Command)
	}
	if !comm.StartCalled {
		t.Fatalf("Expected '%s' to be called, but no remote call was made", expectedCommand)
	}
}

func TestProvisioner_installPackageNonZeroExitCode(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	comm.StartExitStatus = 2
	config["install_modules"] = map[string]string{
		"SomeModuleName":  "1.0.0",
		"SomeModuleName2": "2.0.0",
	}
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = p.installPackage(ctx, ui, comm, "SomeModuleName", "1.0.0")
	if err == nil {
		t.Fatalf("Expected error but got none")
	}
}

func TestProvisionerPrepare_installPackageManagement(t *testing.T) {
	config := testConfig()
	config["install_package_management"] = true
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisioner_installPackageManagement(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	config["install_package_management"] = true
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	err = p.installPackageManagement(ctx, ui, comm)
	if err != nil {
		t.Fatalf("Err: %s", err)
	}
}

func TestProvisioner_installPackageManagementFail(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	config["install_package_management"] = true
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	comm.StartExitStatus = 2
	err = p.installPackageManagement(ctx, ui, comm)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}
}

func TestProvisionerProvision_mofFile(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	mofPath, _ := ioutil.TempDir("/tmp", "packer")
	defer os.Remove(mofPath)

	config["configuration_name"] = "SomeProjectName"
	config["mof_path"] = mofPath
	config["module_paths"] = []string{"."}

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ctx, ui, comm, config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	s := comm.StartCmd.Command
	re := regexp.MustCompile(`powershell \"\& \{ ([a-zA-Z0-9-\/]+).*`)
	command := re.FindStringSubmatch(s)[1]

	bytes, err := ioutil.ReadFile(command)
	if err != nil {
		t.Fatalf(err.Error())
	}
	scriptContents := strings.TrimSpace(string(bytes))

	expectedCommand := `
#
# DSC Runner.
#
# Bootstraps the DSC environment, sets up configuration data
# and runs the DSC Configuration.
#
#
# Set the local PowerShell Module environment path

$absoluteModulePaths = [string]::Join(";", ("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { $_ | Resolve-Path }))
echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })


$script = $("/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" | Resolve-Path)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"


$StagingPath = "/tmp/packer-dsc-pull/mof"


# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`

	if scriptContents != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), scriptContents)
	}
}

func TestProvisionerProvision_noConfigurationParams(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	config["configuration_name"] = "SomeProjectName"
	config["module_paths"] = []string{"."}
	delete(config, "configuration_file")
	delete(config, "configuration_params")

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ctx, ui, comm, config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	s := comm.StartCmd.Command
	re := regexp.MustCompile(`powershell \"\& \{ ([a-zA-Z0-9-\/]+).*`)
	command := re.FindStringSubmatch(s)[1]

	bytes, err := ioutil.ReadFile(command)
	if err != nil {
		t.Fatalf(err.Error())
	}
	scriptContents := strings.TrimSpace(string(bytes))

	expectedCommand := `
#
# DSC Runner.
#
# Bootstraps the DSC environment, sets up configuration data
# and runs the DSC Configuration.
#
#
# Set the local PowerShell Module environment path

$absoluteModulePaths = [string]::Join(";", ("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { $_ | Resolve-Path }))
echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })


$script = $("/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" | Resolve-Path)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"


# Generate the MOF file, only if a MOF path not already provided.
# Import the Manifest
. $script

cd "/tmp/packer-dsc-pull"
$StagingPath = $(Join-Path "/tmp/packer-dsc-pull" "staging")

SomeProjectName -OutputPath $StagingPath 


# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`

	if scriptContents != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), scriptContents)
	}
}

func TestProvisionerProvision_configurationParams(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	configurationParams := map[string]string{
		"-Website": "Beanstalk",
	}
	config["configuration_name"] = "SomeProjectName"
	delete(config, "configuration_file")
	config["configuration_params"] = configurationParams
	config["module_paths"] = []string{"."}

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ctx, ui, comm, config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	s := comm.StartCmd.Command
	re := regexp.MustCompile(`powershell \"\& \{ ([a-zA-Z0-9-\/]+).*`)
	command := re.FindStringSubmatch(s)[1]

	bytes, err := ioutil.ReadFile(command)
	if err != nil {
		t.Fatalf(err.Error())
	}
	scriptContents := strings.TrimSpace(string(bytes))

	expectedCommand := `
#
# DSC Runner.
#
# Bootstraps the DSC environment, sets up configuration data
# and runs the DSC Configuration.
#
#
# Set the local PowerShell Module environment path

$absoluteModulePaths = [string]::Join(";", ("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { $_ | Resolve-Path }))
echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("/tmp/packer-dsc-pull/module-0".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })


$script = $("/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" | Resolve-Path)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"


# Generate the MOF file, only if a MOF path not already provided.
# Import the Manifest
. $script

cd "/tmp/packer-dsc-pull"
$StagingPath = $(Join-Path "/tmp/packer-dsc-pull" "staging")

SomeProjectName -OutputPath $StagingPath -Website "Beanstalk"


# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`

	if scriptContents != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), scriptContents)
	}
}

func TestProvisioner_removeDir(t *testing.T) {
	ctx := context.Background()
	config := testConfig()
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ui := &packer.BasicUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	err = p.removeDir(ctx, ui, comm, "somedir")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	comm.StartExitStatus = 1
	err = p.removeDir(ctx, ui, comm, "somedir")
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

}
