<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    This script will call backup-wsl-bootstrap.ps1 with elevated privileges.

.DESCRIPTION 
    A long description of how the script works and how to use it.
    This step is simply done to inform the user that we need elevation to continue with the bootstrap process.
 
.NOTES 

.COMPONENT 
    Information about PowerShell Modules to be required.
    Tested with PowerShell $PSVersionTable 7.4.1. Make sure to install PowerShell 7.4.1 or later. https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows

.LINK 
    Useful Link to ressources or others.
 
.Parameter ParameterName 
    Description for a parameter in param definition section. Each parameter requires a separate description. The name in the description and the parameter section must match. 

#>
Set-Culture en-US

# Set the log file
#$logfile = "$env:OneDrive\wsl-scripts\"+[Environment]::MachineName+"-get-user-input-bootstrap.log"
$logfile = "$PSScriptRoot\"+[Environment]::MachineName+"-get-user-input-bootstrap.log"

# Start the transcript
Start-Transcript -Path $logfile

# Install the BurntToast module without confirmation
Install-Module -Name BurntToast -Scope CurrentUser -Force

# Import the BurntToast module
Import-Module -Name BurntToast

# User info and input
Add-Type -AssemblyName System.Windows.Forms
$msgBoxInput = [System.Windows.Forms.MessageBox]::Show('This setup needs elevation with Admin By Request' + [Environment]::NewLine + [Environment]::NewLine + 'Click OK to install and click Yes in the Admin by request dialogue that will appear after you click OK','ITM WSL Backup bootstrap utility','OKCancel','Information','Button1','ServiceNotification')
switch ($msgBoxInput) {

    'OK' {    
        # Start-Process pwsh  -Verb RunAs -ArgumentList "-NoExit -File backup-wsl-bootstrap.ps1"
        Start-Process -FilePath "pwsh" -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File backup-wsl-bootstrap.ps1" -Verb RunAs
    }
    'Cancel' {
        # do nothing
    }
}

# Stop the transcript  
exit
Stop-Transcript

