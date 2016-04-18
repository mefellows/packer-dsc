require 'formula'

class PackerProvisionerDsc < Formula
  homepage "https://github.com/mefellows/packer-dsc/"
  version "0.0.1-pre-release"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_amd64.zip"
    sha1 '99d3a857486770821f71b22482d3df56715cefc0'
  else
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_386.zip"
    sha1 '9cf1e701620e355ffd6fa6dd93a9314cfca5321a'
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
