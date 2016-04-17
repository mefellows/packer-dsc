function Get-TargetResource {

  [CmdletBinding()]
  [OutputType([Hashtable])]
  param (
    [parameter(Mandatory)]
    [String]$Subject,

    [parameter(Mandatory)]
    [ValidateSet('CurrentUser','LocalMachine')]
    [String]$StoreLocation,

    [parameter(Mandatory)]
    [ValidateSet('TrustedPublisher','ClientAuthIssuer','Root','CA','AuthRoot','TrustedPeople','My','SmartCardRoot','Trust','Disallowed')]
    [String]$StoreName
  )

  $certInfo = Get-ChildItem -Path "Cert:\$StoreLocation\$StoreName" | Where-Object Subject -eq $Subject

  if ($certInfo) {
    Write-Verbose -Message "The certificate with subject: $Subject exists."
    $ensureResult = 'Present'
  } else {
    Write-Verbose -Message "The certificate with subject: $Subject does not exist."
    $ensureResult = 'Absent'
  }

  $returnValue = @{
    Subject = $certInfo.Subject
    Ensure = $ensureResult
  }

  $returnValue
}

function Set-TargetResource {
  [CmdletBinding()]
  param (
    [parameter(Mandatory)]
    [String]$Subject,

    [parameter(Mandatory)]
    [ValidateSet('CurrentUser','LocalMachine')]
    [String]$StoreLocation,

    [parameter(Mandatory)]
    [ValidateSet('TrustedPublisher','ClientAuthIssuer','Root','CA','AuthRoot','TrustedPeople','My','SmartCardRoot','Trust','Disallowed')]
    [String]$StoreName,

    [parameter(Mandatory)]
    [ValidateSet('Absent','Present')]
    [String]$Ensure

  )

  if ($Ensure -eq 'Present') {

    # TODO: Vendor binary to avoid any runtime issues if that file disappears
    if (-not(Get-Module -Name MrCertificate -ListAvailable)) {
      $Uri = 'https://gallery.technet.microsoft.com/scriptcenter/Self-signed-certificate-5920a7c6/file/101251/1/New-SelfSignedCertificateEx.zip'
      $ModulePath = "$env:ProgramFiles\WindowsPowerShell\Modules\MrCertificate"
      $OutFile = "$env:ProgramFiles\WindowsPowerShell\Modules\MrCertificate\New-SelfSignedCertificateEx.zip"
      Write-Verbose -Message 'Required module MrCertificate does not exist and will now be installed.'
      New-Item -Path $ModulePath -ItemType Directory
      Write-Verbose -Message 'Downloading the New-SelfSignedCertificateEx.zip file from the TechNet script repository.'
      Invoke-WebRequest -Uri $Uri -OutFile $OutFile
      Unblock-File -Path $OutFile
      Write-Verbose -Message 'Extracting the New-SelfSignedCertificateEx.zip file to the MrCertificate module folder.'
      Add-Type -AssemblyName System.IO.Compression.FileSystem
      [System.IO.Compression.ZipFile]::ExtractToDirectory($OutFile, $ModulePath)
      Write-Verbose -Message 'Creating the mrcertificate.psm1 file and adding the necessary content to it.'
      New-Item -Path $ModulePath -Name mrcertificate.psm1 -ItemType File |
      Add-Content -Value '. "$env:ProgramFiles\WindowsPowerShell\Modules\MrCertificate\New-SelfSignedCertificateEx.ps1"'
      Remove-Item -Path $OutFile -Force
    }

    Import-Module -Name MrCertificate
    Write-Verbose -Message "Creating certificate with subject: $Subject"
    New-SelfSignedCertificateEx -Subject $Subject -StoreLocation $StoreLocation -StoreName $StoreName
  } elseif ($Ensure -eq 'Absent') {
    Write-Verbose -Message "Removing the certificate with subject $Subject."
    Get-ChildItem -Path "Cert:\$StoreLocation\$StoreName" |
    Where-Object Subject -eq $Subject |
    Remove-Item -Force
  }
}

function Test-TargetResource {
  [CmdletBinding()]
  [OutputType([Boolean])]
  param (
    [parameter(Mandatory)]
    [String]$Subject,

    [parameter(Mandatory)]
    [ValidateSet('CurrentUser','LocalMachine')]
    [String]$StoreLocation,

    [parameter(Mandatory)]
    [ValidateSet('TrustedPublisher','ClientAuthIssuer','Root','CA','AuthRoot','TrustedPeople','My','SmartCardRoot','Trust','Disallowed')]
    [String]$StoreName,

    [parameter(Mandatory)]
    [ValidateSet('Absent','Present')]
    [String]$Ensure
  )

  $certInfo = Get-ChildItem -Path "Cert:\$StoreLocation\$StoreName" |
  Where-Object Subject -eq $Subject
  switch ($Ensure) {
    'Absent' {$DesiredSetting = $false}
    'Present' {$DesiredSetting = $true}
  }

  if (($certInfo -and $DesiredSetting) -or (-not($certInfo) -and -not($DesiredSetting))) {
    Write-Verbose -Message 'The certificate matches the desired state.'
    [Boolean]$result = $true
  } else {
    Write-Verbose -Message 'The certificate does not match the desired state.'
    [Boolean]$result = $false
  }

  $result
}

Export-ModuleMember -Function *-TargetResource
