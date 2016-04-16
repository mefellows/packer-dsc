// This package implements a provisioner for Packer that executes
// DSC on the remote machine, configured to apply a local manifest
// versus connecting to a DSC push server.
//
// NOTE: This has only been tested on Windows environments
package dsc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Provisioner struct {
	config Config
}

type ExecuteTemplate struct {
	WorkingDir            string
	ConfigurationParams   string
	ConfigurationFilePath string
	ConfigurationName     string
	ModulePath            string
	ManifestFile          string
	ManifestDir           string
	MofPath               string
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"execute_command",
			},
		},
	}, raws...)
	if err != nil {
		return err
	}

	// Set some defaults
	if p.config.ExecuteCommand == "" {
		p.config.ExecuteCommand = `
#
# DSC Runner.
#
# Bootstraps the DSC environment, sets up configuration data
# and runs the DSC Configuration.
#
#
# Set the local PowerShell Module environment path
{{if ne .ModulePath ""}}
$absoluteModulePaths = [string]::Join(";", ("{{.ModulePath}}".Split(";") | ForEach-Object { $_ | Resolve-Path }))
echo "Adding to path: $absoluteModulePaths"
$env:PSModulePath="$absoluteModulePaths;${env:PSModulePath}"
("{{.ModulePath}}".Split(";") | ForEach-Object { gci -Recurse  $_ | ForEach-Object { Unblock-File  $_.FullName} })
{{end}}

$script = $(Join-Path "" "{{.ManifestFile}}" -Resolve)
echo "PSModulePath Configured: ${env:PSModulePath}"
echo "Running Configuration file: ${script}"

{{if eq .MofPath ""}}
# Generate the MOF file, only if a MOF path not already provided.
# Import the Manifest
. $script

cd "{{.WorkingDir}}"
$StagingPath = $(Join-Path "{{.WorkingDir}}" "staging")
{{if ne .ConfigurationFilePath ""}}
$Config = $(iex (Get-Content (Join-Path "{{.WorkingDir}}" "{{.ConfigurationFilePath}}" -Resolve) | Out-String))
{{end}}
{{.ConfigurationName}} -OutputPath $StagingPath {{.ConfigurationParams}}{{if ne .ConfigurationFilePath ""}} -ConfigurationData $Config{{end}}
{{else}}
$StagingPath = "{{.MofPath}}"
{{end}}

# Start a DSC Configuration run
Start-DscConfiguration -Force -Wait -Verbose -Path $StagingPath`
	}

	if p.config.StagingDir == "" {
		p.config.StagingDir = "/tmp/packer-dsc-pull"
	}

	if p.config.WorkingDir == "" {
		p.config.WorkingDir = p.config.StagingDir
	}

	if p.config.ConfigurationParams == nil {
		p.config.ConfigurationParams = make(map[string]string)
	}

	// Validation
	var errs *packer.MultiError
	if p.config.ConfigurationFilePath != "" {
		info, err := os.Stat(p.config.ConfigurationFilePath)
		if err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("configuration_file_path is invalid: %s", err))
		} else if info.IsDir() {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("configuration_file_path must point to a file"))
		}
	}

	if p.config.ManifestDir != "" {
		info, err := os.Stat(p.config.ManifestDir)
		if err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("manifest_dir is invalid: %s", err))
		} else if !info.IsDir() {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("manifest_dir must point to a directory"))
		}
	}

	if p.config.ManifestFile == "" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("A manifest_file must be specified."))
	} else {
		_, err := os.Stat(p.config.ManifestFile)
		if err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("manifest_file is invalid: %s", err))
		}
	}

	for i, path := range p.config.ModulePaths {
		info, err := os.Stat(path)
		if err != nil {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("module_path[%d] is invalid: %s", i, err))
		} else if !info.IsDir() {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("module_path[%d] must point to a directory", i))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	ui.Say("Provisioning with DSC...")
	ui.Message("Creating DSC staging directory...")
	if err := p.createDir(ui, comm, p.config.StagingDir); err != nil {
		return fmt.Errorf("Error creating staging directory: %s", err)
	}

	// Upload configuration_params config if set
	remoteConfigurationFilePath := ""
	if p.config.ConfigurationFilePath != "" {
		var err error
		remoteConfigurationFilePath, err = p.uploadConfigurationFile(ui, comm)
		if err != nil {
			return fmt.Errorf("Error uploading configuration_params config: %s", err)
		}
	}

	// Upload manifest dir if set
	remoteManifestDir := ""
	if p.config.ManifestDir != "" {
		ui.Message(fmt.Sprintf(
			"Uploading manifest directory from: %s", p.config.ManifestDir))
		remoteManifestDir = fmt.Sprintf("%s/manifest", p.config.StagingDir)
		err := p.uploadDirectory(ui, comm, remoteManifestDir, p.config.ManifestDir)
		if err != nil {
			return fmt.Errorf("Error uploading manifest dir: %s", err)
		}
	}

	// Upload all modules
	modulePaths := make([]string, 0, len(p.config.ModulePaths))
	for i, path := range p.config.ModulePaths {
		ui.Message(fmt.Sprintf("Uploading local modules from: %s", path))
		targetPath := fmt.Sprintf("%s/module-%d", p.config.StagingDir, i)
		if err := p.uploadDirectory(ui, comm, targetPath, path); err != nil {
			return fmt.Errorf("Error uploading modules: %s", err)
		}

		modulePaths = append(modulePaths, targetPath)
	}

	// Upload pre-generated MOF
	remoteMofPath := ""
	if p.config.MofPath != "" {
		ui.Message(fmt.Sprintf("Uploading local MOF path from: %s", p.config.MofPath))
		remoteMofPath = fmt.Sprintf("%s/mof", p.config.StagingDir)
		if err := p.uploadDirectory(ui, comm, remoteMofPath, p.config.MofPath); err != nil {
			return fmt.Errorf("Error uploading MOF: %s", err)
		}
	}

	// Upload manifest
	remoteManifestFile, err := p.uploadManifest(ui, comm)
	if err != nil {
		return fmt.Errorf("Error uploading manifest: %s", err)
	}

	// Compile the configuration variables
	configurationVars := make([]string, 0, len(p.config.ConfigurationParams))
	for k, v := range p.config.ConfigurationParams {
		if v == "" {
			configurationVars = append(configurationVars, fmt.Sprintf(`%s `, k))
		} else {
			configurationVars = append(configurationVars, fmt.Sprintf(`%s "%s"`, k, v))
		}
	}

	// Execute DSC script template
	p.config.ctx.Data = &ExecuteTemplate{
		ConfigurationParams:   strings.Join(configurationVars, " "),
		ConfigurationFilePath: remoteConfigurationFilePath,
		ManifestDir:           remoteManifestDir,
		ManifestFile:          remoteManifestFile,
		ModulePath:            strings.Join(modulePaths, ";"),
		WorkingDir:            p.config.WorkingDir,
		ConfigurationName:     p.config.ConfigurationName,
		MofPath:               remoteMofPath,
	}
	command, err := interpolate.Render(p.config.ExecuteCommand, &p.config.ctx)

	if err != nil {
		return err
	}

	cmd := &packer.RemoteCmd{
		Command: command,
	}

	ui.Message(fmt.Sprintf("Running DSC: %s", command))
	if err := cmd.StartWithUi(comm, ui); err != nil {
		return err
	}

	if cmd.ExitStatus != 0 && cmd.ExitStatus != 2 && !p.config.IgnoreExitCodes {
		return fmt.Errorf("DSC exited with a non-zero exit status: %d", cmd.ExitStatus)
	}

	if p.config.CleanStagingDir {
		if err := p.removeDir(ui, comm, p.config.StagingDir); err != nil {
			return fmt.Errorf("Error removing staging directory: %s", err)
		}
	}
	return nil
}

func (p *Provisioner) Cancel() {
	// Just hard quit. It isn't a big deal if what we're doing keeps
	// running on the other side.
	os.Exit(0)
}

func (p *Provisioner) uploadConfigurationFile(ui packer.Ui, comm packer.Communicator) (string, error) {
	ui.Message("Uploading configuration parameters...")
	f, err := os.Open(p.config.ConfigurationFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	path := fmt.Sprintf("%s/%s", p.config.StagingDir, p.config.ConfigurationFilePath)
	if err := comm.Upload(path, f, nil); err != nil {
		return "", err
	}

	return path, nil
}

func (p *Provisioner) uploadManifest(ui packer.Ui, comm packer.Communicator) (string, error) {
	// Create the remote manifest directory...
	ui.Message("Uploading manifest...")
	remoteManifestDir := fmt.Sprintf("%s/manifest", p.config.StagingDir)
	if err := p.createDir(ui, comm, remoteManifestDir); err != nil {
		return "", fmt.Errorf("Error creating manifest directory: %s", err)
	}

	ui.Message(fmt.Sprintf("Uploading manifest file from: %s", p.config.ManifestFile))

	f, err := os.Open(p.config.ManifestFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	manifestFilename := filepath.Base(p.config.ManifestFile)
	remoteManifestFile := fmt.Sprintf("%s/%s", remoteManifestDir, manifestFilename)
	if err := comm.Upload(remoteManifestFile, f, nil); err != nil {
		return "", err
	}
	return remoteManifestFile, nil
}

func (p *Provisioner) createDir(ui packer.Ui, comm packer.Communicator, dir string) error {
	cmd := &packer.RemoteCmd{
		Command: fmt.Sprintf("mkdir '%s' -Force -ErrorAction SilentlyContinue", dir),
	}

	if err := cmd.StartWithUi(comm, ui); err != nil {
		return err
	}

	if cmd.ExitStatus != 0 {
		return fmt.Errorf("Non-zero exit status.")
	}

	return nil
}

func (p *Provisioner) removeDir(ui packer.Ui, comm packer.Communicator, dir string) error {
	cmd := &packer.RemoteCmd{
		Command: fmt.Sprintf("rm -fr '%s'", dir),
	}

	if err := cmd.StartWithUi(comm, ui); err != nil {
		return err
	}

	if cmd.ExitStatus != 0 {
		return fmt.Errorf("Non-zero exit status.")
	}

	return nil
}

func (p *Provisioner) uploadDirectory(ui packer.Ui, comm packer.Communicator, dst string, src string) error {
	if err := p.createDir(ui, comm, dst); err != nil {
		return err
	}

	// Make sure there is a trailing "/" so that the directory isn't
	// created on the other side.
	if src[len(src)-1] != '/' {
		src = src + "/"
	}

	return comm.UploadDir(dst, src, nil)
}
