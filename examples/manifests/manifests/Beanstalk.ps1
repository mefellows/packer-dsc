Configuration Beanstalk
{
  param (
    [string] $MachineName,
    [string] $WebAppPath = "%SystemDrive%\inetpub\wwwroot",
    [string] $HostName   = "localhost"
  )

  Import-DscResource -Module BaseBeanstalkApp

  Node $MachineName
  {
    WindowsFeature IIS
    {
        Ensure = "Present"
        Name = "Web-Server"
    }

    WindowsFeature IISManagerFeature
    {
        Ensure = "Present"
        Name = "Web-Mgmt-Tools"
    }

    WindowsFeature WebApp
    {
        Ensure = "Present"
        Name = "Web-App-Dev"
		    IncludeAllSubFeature = $True
    }

    BeanstalkWebsite sWebsite
    {
        WebAppPath = $WebAppPath
        Hostname   = $HostName
    }
  }
}
