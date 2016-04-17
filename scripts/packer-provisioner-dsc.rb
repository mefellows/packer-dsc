require 'formula'

class PackerProvisionerDsc < Formula
  homepage "https://github.com/mefellows/packer-dsc/"
  version "0.0.1-pre-release"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_amd64.zip"
    sha1 '7a04b6b73e9b7baabed2b941b478d06f38202582'
  else
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_386.zip"
    sha1 '4bf995cf520dd8c9380bfb258bfc9a970f1bad78'
  end

  depends_on :arch => :intel

  def install
    pluginpath = Pathname.new("#{ENV['HOME']}/.packer.d/plugins")

    unless File.directory?(pluginpath)
      mkdir_p(pluginpath)
    end

    cp_r Dir["*"], pluginpath
    bin.install Dir['*']
  end

  test do
    minimal = testpath/"minimal.json"
    minimal.write <<-EOS.undent
    {
      "builders": [
        {
          "type": "null",
          "ssh_host":     "foo",
      		"ssh_username": "bar",
      		"ssh_password": "baz"
        }
      ],
      "provisioners": [
        {
          "type": "dsc",
          "configuration_name": "Beanstalk",
          "configuration_file": "manifests/BeanstalkConfig.psd1",
          "manifest_file": "manifests/Beanstalk.ps1",
          "module_paths": [
            "modules"
          ],
          "configuration_params": {
            "-WebAppPath": "c:\\tmp",
            "-MachineName": "localhost"
          }
        }
      ]
    }
    EOS
    system "packer", "validate", minimal
  end
end
