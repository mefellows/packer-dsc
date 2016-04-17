require 'formula'

class PackerProvisionerDsc < Formula
  homepage "https://github.com/mefellows/packer-dsc/"
  version "0.0.1-pre-release"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_amd64.zip"
    sha1 'fb347ef854c2020d182f0661b04842d62f6d3cab'
  else
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_386.zip"
    sha1 '22e1b8facb397c51d8fd840806cb571a550adeb4'
  end

  depends_on :arch => :intel

  def install
    pluginpath = Pathname.new("/Users/#{ENV['USER']}/.packer.d/plugins")

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
