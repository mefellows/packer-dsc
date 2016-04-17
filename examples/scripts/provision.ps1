iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1'))
choco install 7zip.commandline -y

# Download the xWebAdministration
# + Install DSC Modules into PS module path
(New-Object System.Net.WebClient).DownloadFile('https://github.com/PowerShell/xWebAdministration/archive/1.9.0.0-PSGallery.zip','c:\xmodules.zip')
7z -y x C:\xmodules.zip -o"${env:ProgramFiles}\WindowsPowerShell\Modules"
mv "${env:ProgramFiles}\WindowsPowerShell\Modules\xWebAdministration-1.9.0.0-PSGallery"  "${env:ProgramFiles}\WindowsPowerShell\Modules\xWebAdministration"
rm c:\xmodules.zip

# Show what resources are now available
Get-DSCResource
