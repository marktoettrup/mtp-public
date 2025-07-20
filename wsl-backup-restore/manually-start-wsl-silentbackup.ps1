<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    

.DESCRIPTION 
    A long description of how the script works and how to use it.

 
.NOTES 
    The script is triggered during bootstrap (setup.exe).

.COMPONENT 
    Information about PowerShell Modules to be required.
    Tested with PowerShell $PSVersionTable 7.4.1. Make sure to install PowerShell 7.4.1 or later. https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows

.LINK 
    Useful Link to ressources or others.
    Based on https://lazyadmin.nl/powershell/how-to-create-a-powershell-scheduled-task/
 
.Parameter ParameterName 
    Description for a parameter in param definition section. Each parameter requires a separate description. The name in the description and the parameter section must match. 
    This script takes no input parameters
#>

Set-Culture en-US

# Set the log file
$logfile = "$env:OneDrive\wsl-scripts\"+[Environment]::MachineName+"-manually-start-wsl-backup.log"

# Start the transcript
Start-Transcript -Path $logfile
try {
# Install the BurntToast module without confirmation
Install-Module -Name BurntToast -Scope CurrentUser -Force

# Import the BurntToast module
Import-Module -Name BurntToast
New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "WSL backup starting..."
# User info and input
Add-Type -AssemblyName System.Windows.Forms

# Task name
$TaskName = "Backup WSL"

# Script paths
# $ScriptPath = "$FilePath\backup-wsl.ps1"

# Check if the task already exists
if (Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue) {
    # If the task exists, unregister (delete) it
    # Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
    Write-Host "The backup WSL Schedule Task exists and will be started"
    Start-ScheduledTask -TaskName $TaskName
} else {
    Write-Host "The backup WSL Schedule Task does not exist"
    # Execute the setup.exe file
    Add-Type -AssemblyName PresentationFramework
    $msgBoxInput = [System.Windows.MessageBox]::Show('The WSL backup solution needs to be re-installed' + [Environment]::NewLine + 'Click Yes to reinstall - this will also allow you to run the backup' + [Environment]::NewLine + 'Click No to close', 'Confirmation', 'YesNo', 'Information')
    
    switch ($msgBoxInput) {
    'Yes' {
        # Run the task
        Start-Process -FilePath "$env:OneDrive\wsl-scripts\setup.exe" -Wait
    }
    'No' {
        # Do nothing
    }
    }
}

# Stop the transcript  
Stop-Transcript


}
catch {
    Write-Host "Error occurred: $_"
    # Handle the error (e.g., log it, display a message, etc.)
    Stop-Transcript
}



