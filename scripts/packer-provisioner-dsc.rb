require 'formula'

class PackerProvisionerDsc < Formula
  homepage "https://github.com/mefellows/packer-dsc/"
  version "0.0.3-pre-release"

  if Hardware::CPU.is_64_bit?
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_amd64.zip"
    sha1 'a9e619b230cdf5b8ff1bb65d56d652d14eaa80e7f83a4fc6bc48fe9c524c89a3'
  else
    url "https://github.com/mefellows/packer-dsc/releases/download/#{version}/darwin_386.zip"
    sha1 '4436241b87377e389d23ed697928dab5e14013294258ff1d22e02d6075ca4e4b'
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
