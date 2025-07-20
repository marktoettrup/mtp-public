<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    Called if the user wants to disable the WSL backup schedule.

.DESCRIPTION 
    A long description of how the script works and how to use it.
    This script will disable the WSL backup schedule.
 
.NOTES 
    The script is triggered from backup-wsl.ps1.
    To enabe, run setup.exe

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
$logfile = "$env:OneDrive\wsl-scripts\"+[Environment]::MachineName+"-backup-wsl-schedule-disable.log"

# Start the transcript
Start-Transcript -Path $logfile

$TaskName = "Backup WSL"
Stop-ScheduledTask -TaskName $TaskName
Disable-ScheduledTask -TaskName $TaskName
Add-Type -AssemblyName PresentationFramework
[System.Windows.MessageBox]::Show('The WSL backup schedule has been permanently disabled', 'Confirmation', 'OK', 'Information')
# Stop the transcript  
Stop-Transcript
