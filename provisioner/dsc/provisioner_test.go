package dsc

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/packer/packer"
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
		"manifest_file": filename,
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

	delete(config, "configuration_file_path")
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test with a good one
	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("error tempfile: %s", err)
	}
	defer os.Remove(tf.Name())

	config["configuration_file_path"] = tf.Name()
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

func TestProvisionerProvision_mofFile(t *testing.T) {
	config := testConfig()
	ui := &packer.MachineReadableUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	mof_path, _ := ioutil.TempDir("/tmp", "packer")
	defer os.Remove(mof_path)

	config["configuration_name"] = "SomeProjectName"
	config["mof_path"] = mof_path
	config["module_paths"] = []string{"."}

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ui, comm)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

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


$script = $(Join-Path "" "/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" -Resolve)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"


$StagingPath = "/tmp/packer-dsc-pull/mof"


# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`

	if strings.TrimSpace(comm.StartCmd.Command) != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), strings.TrimSpace(comm.StartCmd.Command))
	}
}

func TestProvisionerProvision_noConfigurationParams(t *testing.T) {
	config := testConfig()
	ui := &packer.MachineReadableUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	config["configuration_name"] = "SomeProjectName"
	config["module_paths"] = []string{"."}

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ui, comm)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

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


$script = $(Join-Path "" "/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" -Resolve)
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

	if strings.TrimSpace(comm.StartCmd.Command) != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), strings.TrimSpace(comm.StartCmd.Command))
	}
}

func TestProvisionerProvision_configurationParams(t *testing.T) {
	config := testConfig()
	ui := &packer.MachineReadableUi{
		Writer: ioutil.Discard,
	}
	comm := new(packer.MockCommunicator)
	configurationParams := map[string]string{
		"-Website": "Beanstalk",
	}
	config["configuration_name"] = "SomeProjectName"
	config["configuration_params"] = configurationParams
	config["module_paths"] = []string{"."}

	// Test with valid values
	p := new(Provisioner)
	err := p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Provision(ui, comm)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

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


$script = $(Join-Path "" "/tmp/packer-dsc-pull/manifest/packer-dsc-pull-manifest" -Resolve)
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

	if strings.TrimSpace(comm.StartCmd.Command) != strings.TrimSpace(expectedCommand) {
		t.Fatalf("Expected:\n\n%s\n\nbut got: \n\n%s", strings.TrimSpace(expectedCommand), strings.TrimSpace(comm.StartCmd.Command))
	}
}
