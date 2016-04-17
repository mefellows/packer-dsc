Configuration BeanstalkWebsite
{
    param
    (
        [String]$WebAppPath             = "%SystemDrive%\inetpub\wwwroot",
        [String]$WebSiteName            = "Default Web Site",
        [String]$HostNameSuffix         = "dev",
        [String]$HostName               = "beanstalk.${HostNameSuffix}",
        [String]$ApiAppPoolName         = "beanstalk-API",
        [HashTable]$AuthenticationInfo  = @{Anonymous = "true"; Basic = "false"; Digest = "false"; Windows = "false"}
    )

    Import-DscResource -Module xWebAdministration

    # Copy the website content
    # This is only for vagrant testing purposes
    File WebContent
    {
        Ensure          = "Present"
        Contents        = "<h1>Hello from Vagrant</h1>"
        DestinationPath = "c:\tmp\index.html"
        Type            = "File"
    }

    # Stop the default website (beanstalk need this?)
    # Create the new Website with HTTPS
    xWebsite DefaultSite
    {
        Ensure          = "Present"
        Name            = $WebSiteName
        State           = "Started"
        PhysicalPath    = $WebAppPath
        BindingInfo     = @(
            MSFT_xWebBindingInformation
            {
                Protocol              = "HTTP"
                Port                  = 80
            }
        )
        DependsOn = "[File]WebContent"
    }
}

# Get a certificate thumbprint given a subject
# e.g. CN=www.foo.com
function getCertificateThumbprint($subject) {
  (Get-ChildItem -Path Cert:\LocalMachine\My | Where-Object {$_.Subject -match $subject}).Thumbprint;
}
