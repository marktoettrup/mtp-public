
<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    This script will backup and restore WSL distros  

.DESCRIPTION 
    A long description of how the script works and how to use it.
    input param 1 must contain either "backup" or "restore" 
    input param 2 must contain the directory where the backup will be stored or the directory where the backup is located
    the script will check if WSL is present and if so, it will attempt to backup or restore the WSL distros
    If the backup directory does not exist, it will be created.
    At the end, your should see a windows desktop notification with the result, either "WSL backup completed successfully" or "WSL backup error - please check the log at $logfile".
    The script rotates the backup directory, keeping the last three backups.
    By default, the backups are stored in the user's OneDrive folder, ensuring that the backups are included in the OneDrive sync.
 
.NOTES 
   Copy and paste the contents of backup-wsl-schedule.ps1 into an elevated PowerShell window to create a scheduled task to run the backup-wsl.ps1 script at startup.

   Due to wsl.exe outputting unicode to stdout - see https://github.com/microsoft/WSL/issues/4180 - it is required to change the console output encoding to unicode
   It is possible to circumvent the unicode console encoding by using the WSL to list the distros and then uses iconv to convert the output to ascii, ie: $distros = wsl -l --quiet | wsl iconv -c -f utf16 -t ascii
   However, the iconv command is not available in all WSL distros, so the solution is instead to change the console output encoding to unicode and then change it back to the original encoding after the wsl command has been executed
   The wsl command is then executed and the output is split into an array of distros
   The console output encoding is then changed back to the original encoding

.COMPONENT 
    Information about PowerShell Modules to be required.
    Tested with PowerShell $PSVersionTable 7.4.1. Make sure to install PowerShell 7.4.1 or later. https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows
    Dependency to windows task scheduler. 

.LINK 
    Useful Link to ressources or others.
    https://wiki/display/ToolsWSL
    based on https://gist.github.com/mohatb/1177b517f986e6028cc35c0a6c9f8787#file-backup_restore-wsl-bash
    fix based on https://github.com/microsoft/WSL/issues/4180
    fix based on https://github.com/microsoft/WSL/issues/4607
 
.Parameter ParameterName 
    Description for a parameter in param definition section. Each parameter requires a separate description. The name in the description and the parameter section must match. 
    input 1: backup or restore
    input 2: the directory where the backup will be stored or the directory where the backup is located
#>

[CmdletBinding()]
param (
  [Parameter(Mandatory=$true)]
  [string]$Command,
  [Parameter(Mandatory=$true)]
  [string]$BackupDirectory
)

# Set the culture to en-US
Set-Culture en-US

# Set the log file with a timestamp
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$logfile = "$env:OneDrive\wsl-scripts\" + [Environment]::MachineName + "-backup-wsl-$timestamp.log"


# Start the transcript
Start-Transcript -Path $logfile

function Backup-WSL {
  [CmdletBinding()]
  param (
    [Parameter(Mandatory=$true)]
    [string]$BackupDirectory,
    [string]$Command
  )
  # Get the current date and time - to be added to the backup directory name
  $timestamp = Get-Date -Format "yyyy-MM-dd-HH-mmss"
  # Set the parent directory
  $parentDirectory = $BackupDirectory 
  # Add the timestamp to the backup directory
  $BackupDirectory = $BackupDirectory + "\" + [Environment]::MachineName + "\" + $timestamp

  Write-Host "WSL backups will be stored in:" $BackupDirectory

  # Create backup directory if it doesn't exist
  if (!(Test-Path $BackupDirectory)) {
    New-Item -ItemType Directory -Path $BackupDirectory | Out-Null
  }
  # Invoke-Item -Path $BackupDirectory

  # Get list of WSL distros using wsl -l --quiet
  # Change the console output encoding to unicode and then change it back to the original encoding after the wsl command has been executed
  # The wsl command is then executed and the output is split into an array of distros
  # The console output encoding is then changed back to the original encoding
  $console = ([console]::OutputEncoding)
  [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
  $distros = (wsl --list --quiet) -split '\s+'
  [console]::OutputEncoding = $console

  # Create an array of the results returned when exporting wsl distros
  $wsl_export_results = @()

  # Pause OneDrive sync for two hours max if possible to allow for the backup to to be done before resuming after all distros have been exported and 7zip'd
  function Get-OneDriveExePath {
    $pathsToCheck = @(
        "$env:LOCALAPPDATA\Microsoft\OneDrive\OneDrive.exe",           # Per-user install
        "$env:PROGRAMFILES\Microsoft OneDrive\OneDrive.exe",           # System-wide install (your case)
        "$env:PROGRAMFILES(X86)\Microsoft OneDrive\OneDrive.exe"       # 32-bit fallback
    )
    foreach ($path in $pathsToCheck) {
        if (Test-Path $path) {
            return $path
        }
    }
    return $null
    }

    $onedriveExe = Get-OneDriveExePath

    if ($onedriveExe) {
        Write-Host "Pausing OneDrive sync for two hours (default) using: $onedriveExe"
        & $onedriveExe /pause
    } else {
        Write-Warning "OneDrive.exe not found — skipping pause/resume"
    }

  # Export each/all WSL distro
  foreach ($distro in $distros -ne '') {
    Write-Host "starting backup of WSL named:" $distro

    $console = ([console]::OutputEncoding)
    [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding

    New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "starting backup of WSL named: $distro"

    # Export distro to .tar
    $distroArray = (wsl --export "$distro" "$BackupDirectory\$distro.tar") -split '\r\n'
    [console]::OutputEncoding = $console

    # Define 7-Zip path
    $sevenZipPath = "C:\Program Files\7-Zip\7z.exe"

    # Check if 7-Zip is installed
    if (Test-Path $sevenZipPath) {
        Write-Host "7-Zip found. Compressing $distro.tar to $distro.7z"
        
        # Compress the tar file using 7-Zip
        & $sevenZipPath a "$BackupDirectory\$distro.7z" "$BackupDirectory\$distro.tar"

        # If compression succeeded, delete the original tar
        if (Test-Path "$BackupDirectory\$distro.7z") {
            Remove-Item "$BackupDirectory\$distro.tar" -Force
            Write-Host "Deleted original tar file: $distro.tar"
        } else {
            Write-Warning "Compression failed or .7z not created for $distro"
        }
    } else {
        Write-Warning "7-Zip not found at $sevenZipPath – skipping compression for $distro"
    }


    # Collect output (optional logging)
    foreach ($distro_result in $distroArray) {
        Write-Host $distro_result
        $wsl_export_results  += $distro_result
    }
}

  # Resume OneDrive sync
  & $onedriveExe /resume


  # # foreach ($distro in $distros -ne '') {
  # #   Write-Host "starting backup of WSL named:" $distro

  # #   $console = ([console]::OutputEncoding)
  # #   [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
  # #   # Create an array of the results returned when exporting wsl distros
  # #   # New-BurntToastNotification -AppLogo C:\Users\mtp\projects\wsl-backup-restore\tux.ico -Text "Don't forget to smile!",'Your script ran successfully, celebrate!'
  # #   New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "starting backup of WSL named: $distro"
  # #   $distroArray = (wsl --export "$distro" "$BackupDirectory\$distro.tar") -split '\r\n'
  # #   [console]::OutputEncoding = $console
  # #   foreach ($distro_result in $distroArray) {
  # #     Write-Host $distro_result
  # #     # Add the result to the array
  # #     $wsl_export_results  += $distro_result
  # #   }
  # # }    

  # For logging purposes, loop through the array and output each line
  Write-Host "checking all wsl backup results"
  foreach ($result in $wsl_export_results) {
    Write-Host $result
  }

  $allSuccessful = $wsl_export_results -match "success"
  if ($allSuccessful) {
    # Backup directory rotation
    # Get all child directories
    $parentDirectory = $parentDirectory + "\" + [Environment]::MachineName
    $childDirectories = Get-ChildItem -Path $parentDirectory -Directory
    # Check if there are more than five child directories
    if ($childDirectories.Count -gt 5) {
      # Sort child directories by creation time (oldest first)
      $sortedDirectories = $childDirectories | Sort-Object CreationTime
      # Determine the number of directories to delete
      $directoriesToDelete = $sortedDirectories.Count - 3
      # Delete the oldest directories
      for ($i = 0; $i -lt $directoriesToDelete; $i++) {
        $directoryToDelete = $sortedDirectories[$i]
        Remove-Item -Path $directoryToDelete.FullName -Recurse -Force
        Write-Host "Deleted directory: $($directoryToDelete.FullName)"
      }
    } else {
      Write-Host "No action needed. There are $($childDirectories.Count) child directories."
    }

    # Check the export status of all distros and report the result to the user using BurntToast notifications
    Write-Host "WSL backup completed successfully"

    New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "WSL backup completed successfully"
    # # Allow the user to open explorer to see the backup directory
    # Add-Type -AssemblyName PresentationFramework
    # $msgBoxInput = [System.Windows.MessageBox]::Show('The WSL backup has completed on this Windows host (' + [Environment]::MachineName + ')' + [Environment]::NewLine + [Environment]::NewLine + 'Click Yes to open explorer to see the backups' + [Environment]::NewLine + 'Click No to close', 'Confirmation', 'YesNo', 'Information')
 
    # switch ($msgBoxInput) {
    #  'Yes' {
    #      # Open explorer to the backup directory
    #      Invoke-Item -Path $BackupDirectory
    #  }
    #  'No' {
    #      # Do nothing
    #  }
    # }

    # Notify the user (optional)
    New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "WSL backup completed successfully"

    # Open the backup folder in File Explorer
    Invoke-Item -Path $BackupDirectory

    # Start all WSL distros
    Write-Host "Attempting to start all WSL distros after backup..."

    $console = ([console]::OutputEncoding)
    [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
    $distros = (wsl --list --quiet) -split '\s+'
    [console]::OutputEncoding = $console

    foreach ($distro in $distros -ne '') {
        try {
            Write-Host "Starting WSL distro: $distro"
            $console = ([console]::OutputEncoding)
            [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding            
            wsl -d $distro
            [console]::OutputEncoding = $console            
        } catch {
            Write-Warning "Failed to start WSL distro: $distro. Error: $_"
        }
    }


  } else {
    Write-Host "WSL backup error - please check the log at $logfile"
    New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "WSL backup error - please check the log at $logfile"
  }
}

#Restore WSL
function Restore-WSL {
  [CmdletBinding()]
  param (
    [Parameter(Mandatory=$true)]
    [string]$BackupDirectory
  )

  # Get list of backup files
  $backupFiles = Get-ChildItem $BackupDirectory

  # Restore each backup file
  foreach ($backupFile in $backupFiles) {
    $distro = $backupFile.BaseName -replace '.tar',''
    wsl --import $distro "$BackupDirectory\$backupFile"
  }
}

# Main function
function Main {
  #[CmdletBinding()]
  # param (
  #   [Parameter(Mandatory=$true)]
  #   [string]$Command,
  #   [Parameter(Mandatory=$true)]
  #   [string]$BackupDirectory
  # )
  if ($Command -match 'backup') {
    # Call the backup function
    Backup-WSL -BackupDirectory $BackupDirectory -Command $Command
  }
  elseif ($Command -eq 'restore') {
    # Call the restore function
    Restore-WSL -BackupDirectory $BackupDirectory
  }
  else {
    Write-Error 'Error: Invalid command'
    return 1
  }
}

# Install the BurntToast module without confirmation
Install-Module -Name BurntToast -Scope CurrentUser -Force

# Import the BurntToast module
Import-Module -Name BurntToast


# check for WSL status on the system
# expect to see "Usage: wsl.exe [Argument]" if WSL is present
$wsl_status = $false

$console = ([console]::OutputEncoding)
[console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
$wslstate = (wsl --help) -split '\r\n'

foreach ($item in $wslstate) {
  if ($item.Contains("Usage: wsl.exe [Argument]")) {
    $wsl_status = $true
  } else {
     # do nothing
  }
}
[console]::OutputEncoding = $console

if ($wsl_status) {
  Write-Host "WSL is present and backup will be attempted" 
  # get the number of WSL installed on the system
  $console = ([console]::OutputEncoding)
  [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
  $distro_counts = (wsl --list --quiet) -split '\s+'
  [console]::OutputEncoding = $console
  $wslCount = $distro_counts.Count
  # Convert array to comma-separated string
  $commaSeparatedWSLNames = $distro_counts -join ","

  # get the trigger of the scheduled task named "Backup WSL"
  $TaskName = "Backup WSL"
  $task = Get-ScheduledTask -TaskName $TaskName
  Write-Host $task
  $taskTrigger = $task.Triggers[0]  # Assuming there's only one trigger
  Write-Host $taskTrigger
  $trigger_delay = $taskTrigger.delay
  if ($trigger_delay) {
    $taskTrigger = "the trigger is at $taskTrigger with a delay of " + $trigger_delay
  } else {
    $taskTrigger = "the trigger is at $taskTrigger"
  }

  # Ask the user if the backup should be perform now, later or never
  # [string]$Command,
  # [string]$BackupDirectory
  Write-Host "Command contains " $Command
  if ($Command -eq 'backup') {
  Add-Type -AssemblyName System.Windows.Forms
  $msgBoxInput = [System.Windows.Forms.MessageBox]::Show('You have ' + $wslCount + ' WSL system(s) installed on this Windows host (' + [Environment]::MachineName + ')' + [Environment]::NewLine + 'Your WSL system(s): ' + $commaSeparatedWSLNames + [Environment]::NewLine + [Environment]::NewLine + 'Click Yes to perform backup of your WSL system(s) now' + [Environment]::NewLine + 'Your WSL system(s) will be closed during the process and backups will be placed in the wsl-backup folder in your One-Drive ' + [Environment]::NewLine + [Environment]::NewLine   + 'Click No, if you want to skip this backup. WSL backup will be attempted again later'  + [Environment]::NewLine + 'According to the trigger of the Scheduled Task named: ' + $TaskName + " " + $taskTrigger + [Environment]::NewLine + [Environment]::NewLine + 'Click Cancel if you want to disable this WSL backup schedule permanently. No further WSL backup attempts will be made','ITM WSL Backup utility','YesNoCancel','Information','Button1','ServiceNotification')

    switch ($msgBoxInput) {

        'Yes' {
          New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "WSL backup starting"
          Main -Command $Command -BackupDirectory $BackupDirectory
        }
        'No' {
          New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "The WSL backup process has been closed"
        }
        'Cancel' {
          if (Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue) {
            # If the task exists, disable it
            Start-Process pwsh -Verb RunAs -ArgumentList "-NoExit -File $env:OneDrive\wsl-scripts\backup-wsl-schedule-disable.ps1"
            New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "The WSL backup schedule has been permanently disabled. Run setup.exe if you want to enable it again"
        }
        }

    }
    } elseif ($Command -eq 'silentbackup') {
      Write-Host "silent backup triggered"
      New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "Silent WSL backup starting"
      Main -Command $Command -BackupDirectory $BackupDirectory
    } 
    else {
      Write-Host "Command is not backup or silentbackup"
    }

  } else {
    Write-Host "WSL not present and no further actions are taken"
    # return 1
  }


# # Check if WSL feature is enabled. If so, continue. oh no, this check requires elevation and I can't seem to get that to run along side the Admin by Request. I've create as cherwell to security
# $wslfeature = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux
# Write-Host "WSL state = " $wslfeature.State 
#  if ($wslfeature.State -eq 'Enabled') {
#   Write-Host "WSL is present and we will do backup" 
#    Main -Command $Command -BackupDirectory $args[1]
#  }
# else {
#   Write-Host "WSL not present and no furhter actions are taken"
#   #return 1
#  }



Stop-Transcript
