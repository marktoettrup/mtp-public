<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    This script is triggered from start-here-bootstrap.ps1 and it will prepare the bootstrap for the backup and restore WSL solution 

.DESCRIPTION 
    A long description of how the script works and how to use it.
    If powershell 7 is not present on the system, the script will install it. The script will start the bootstrap script.

.COMPONENT 
    Information about PowerShell Modules to be required.
    Tested with PowerShell $PSVersionTable 7.4.1. Make sure to install PowerShell 7.4.1 or later. https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows
#>
Set-Culture en-US

# Set the log file
#$logfile = "$env:OneDrive\wsl-scripts\"+[Environment]::MachineName+"-setup.log"
$logfile = "$PSScriptRoot\"+[Environment]::MachineName+"-setup.log"
Start-Transcript -Path $logfile

# Check if PowerShell 7 is installed
if (!(Test-Path "$env:ProgramFiles\\PowerShell\\7")) {
  $confirmwslnamecust = Read-Host "Powershell 7 was not found on the system. Please type y to confirm installation and proceed or n to abort"
    if ($confirmwslnamecust -eq 'y' -or $confirmwslnamecust -eq 'Y') {
      Write-host "OK, proceeding to install PowerShell version 7...'" -ForegroundColor Green
      Start-Process -FilePath "powershell" -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File install-pwsh7.ps1" -Verb RunAs -Wait
    } else {
      Write-host "Aborting..." -ForegroundColor Red
      exit
    }

    } else {
    # do nothing
}

# start the bootstrap script
Start-Process pwsh -ArgumentList "-File get-user-input-bootstrap.ps1"
Stop-Transcript
exit