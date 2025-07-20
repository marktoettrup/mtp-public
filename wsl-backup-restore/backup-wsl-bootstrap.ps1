<# 
.SYNOPSIS
    A summary of how the script works and how to use it.
    This script is triggered by get-user-input-bootstrap.ps1 and it creates necessary folders and windows task in the Scheduler to run WSL backup at startup

.DESCRIPTION 
    A long description of how the script works and how to use it.
    If a task with the name "Backup WSL" already exists, it will be unregistered (deleted) and a new task will be created.
    You can change the task trigger to run the script at a different time by changing the $trigger variable.
    This script creates the necessary folders and copies the backup-wsl.ps1 script to the $env:OneDrive\wsl-scripts\ folder.
 
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
#$logfile = "$env:OneDrive\wsl-scripts\"+[Environment]::MachineName+"-backup-wsl-bootstrap.log"
$logfile = "$PSScriptRoot\"+[Environment]::MachineName+"-backup-wsl-bootstrap.log"

# Start the transcript
Start-Transcript -Path $logfile
try {

# get task status and last run result    
function Get-TaskStatus {
    param (
        [string]$TaskName
    )

    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue

    if ($task) {
        $info = Get-ScheduledTaskInfo -TaskName $TaskName
        $errorMessage = [ComponentModel.Win32Exception]::new($info.LastTaskResult).Message

        # Assume only one action (usually true for these types of tasks)
        $action = $task.Actions[0]

        [PSCustomObject]@{
            TaskName        = $TaskName
            State           = $info.State
            LastRunTime     = $info.LastRunTime
            LastTaskResult  = "$($info.LastTaskResult) - $errorMessage"
            NextRunTime     = $info.NextRunTime
            Execute         = $action.Execute
            WorkingDirectory = $action.WorkingDirectory
            Arguments       = $action.Arguments
        }
    } else {
        Write-Warning "Scheduled task '$TaskName' not found."
    }
}


# Install the BurntToast module without confirmation
Install-Module -Name BurntToast -Scope CurrentUser -Force

# Import the BurntToast module
Import-Module -Name BurntToast

# User info and input
Add-Type -AssemblyName System.Windows.Forms

$FilePath =  "$env:OneDrive\wsl-scripts"
if (!(Test-Path $FilePath)) {
New-Item -ItemType Directory -Path $FilePath | Out-Null
}

$BackupDir = "$env:OneDrive\wsl-backup"
if (!(Test-Path $BackupDir)) {
New-Item -ItemType Directory -Path $BackupDir | Out-Null
} 

$sourceFolder = "$PSScriptRoot"

Write-Host "Copy the content of the source folder $sourceFolder to the destination $FilePath folder (including subdirectories)"
# delete all files and folder in the destination folder
Remove-Item -Path $FilePath\* -Force -Recurse
# copy the content of the source folder to the destination folder except .log files and .git folder and .gitignore file
Copy-Item -Exclude "*.log", ".git", ".gitignore" -Path $sourceFolder\* -Destination $FilePath -Force -Recurse

################### scheduled task ######################

# Create a scheduled task with two triggers:
# Set the execution policy
# Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force

# Task name
$TaskName = "Backup WSL"

# Script paths
$ScriptPath = "$FilePath\backup-wsl.ps1"

# Check if the task already exists
if (Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue) {
    # If the task exists, unregister (delete) it
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
}

# Define the action that the task will execute
#$Action = New-ScheduledTaskAction -Execute "pwsh.exe" -Argument "-WindowStyle Hidden -ExecutionPolicy Bypass -File $ScriptPath backup $BackupDir" -WorkingDirectory $BackupDir
#$Action = New-ScheduledTaskAction -Execute "pwsh.exe" -Argument "-WindowStyle Hidden -ExecutionPolicy Bypass -File `"$ScriptPath`" backup `"$BackupDir`"" -WorkingDirectory `"$BackupDir`"
# Construct the command as a string (quote everything properly)
$pwshPath = "C:\Program Files\PowerShell\7\pwsh.exe"

$quotedScriptPath = '"{0}"' -f $ScriptPath
$quotedBackupDir = '"{0}"' -f $BackupDir
$argString = "-WindowStyle Hidden -ExecutionPolicy Bypass -File $quotedScriptPath backup $quotedBackupDir"

$Action = New-ScheduledTaskAction -Execute $pwshPath -Argument $argString -WorkingDirectory $BackupDir

# Define the logon trigger
# $LogonTrigger = New-ScheduledTaskTrigger -AtLogOn
# $LogonTrigger.DaysOfWeek = 'Tuesday', 'Friday'  # Set the trigger to run on Tuesdays and Fridays

# Add a five-minute delay
# $LogonTrigger.Delay = "PT5M"  # 5 minutes delay

# Define the daily trigger at 10:00 AM
# $DailyTrigger = New-ScheduledTaskTrigger -Daily
# $DailyTrigger.At = "10:00 AM"
$DailyTrigger = New-ScheduledTaskTrigger -Weekly -WeeksInterval 2 -DaysOfWeek Monday -At 10am

# Combine both triggers
# $CombinedTriggers = @($LogonTrigger, $DailyTrigger)

# Define the principal for the task (the user account under which the task runs)
$Principal = New-ScheduledTaskPrincipal -UserId "$env:USERDOMAIN\$env:USERNAME" -LogonType Interactive -RunLevel Highest
#$Principal = New-ScheduledTaskPrincipal -UserId "nt authority\system" -LogonType Interactive -RunLevel Highest

# Register the scheduled task
Register-ScheduledTask -Action $Action -Trigger $DailyTrigger -TaskName $TaskName -Description "Task to run WSL backup every second Monday at 10 AM" -Principal $Principal

################### scheduled task created ######################

################### scheduled task SILENT ######################

# Task name for silent backup
$TaskNameSilent = "Backup WSL Silent"

# Script path (already defined earlier, but we'll reuse to be safe)
$ScriptPath = "$FilePath\backup-wsl.ps1"

# Full path to pwsh.exe
$pwshPath = "C:\Program Files\PowerShell\7\pwsh.exe"

# Ensure the silent backup task doesn't already exist
if (Get-ScheduledTask -TaskName $TaskNameSilent -ErrorAction SilentlyContinue) {
    Unregister-ScheduledTask -TaskName $TaskNameSilent -Confirm:$false
}

# Properly quote paths for arguments
$quotedScriptPathSilent = '"{0}"' -f $ScriptPath
$quotedBackupDirSilent = '"{0}"' -f $BackupDir
$argStringSilent = "-WindowStyle Hidden -ExecutionPolicy Bypass -File $quotedScriptPathSilent silentbackup $quotedBackupDirSilent"

# Create the action
$ActionSilent = New-ScheduledTaskAction -Execute $pwshPath -Argument $argStringSilent -WorkingDirectory $BackupDir

# Define the weekly trigger for every Sunday at 10:00 AM
$WeeklyTriggerSilent = New-ScheduledTaskTrigger -Weekly -DaysOfWeek Sunday -At 10am

# Define the principal for the silent backup task
$PrincipalSilent = New-ScheduledTaskPrincipal -UserId "$env:USERDOMAIN\$env:USERNAME" -LogonType Interactive -RunLevel Highest

# Register the scheduled task for silent backup
Register-ScheduledTask -Action $ActionSilent -Trigger $WeeklyTriggerSilent -TaskName $TaskNameSilent -Description "Task to run WSL silent backup every Sunday at 10 AM" -Principal $PrincipalSilent

################### scheduled task SILENT created ######################

# nootify the user of the installation
New-BurntToastNotification -AppLogo $env:OneDrive\wsl-scripts\tux.ico -Text "The ITM WSL Backup has been installed and scheduled to run every second Monday at 10 AM, and silently, every Sunday at 10.00 AM"

################### DESKTOP SHORTCUT ######################

# Define the source file location
$SourceFilePath = "$FilePath\manually-start-wsl-backup.exe"

# Define the shortcut file location and name
$ShortcutPath = "$FilePath\manually-start-wsl-backup.lnk"

# Create a new WScript.Shell object
$WScriptObj = New-Object -ComObject "WScript.Shell"

# Create the shortcut using the specified path
$shortcut = $WScriptObj.CreateShortcut($ShortcutPath)

# Set the target path for the shortcut
$shortcut.TargetPath = $SourceFilePath

# Save the shortcut
$shortcut.Save()
$desktop_dest = [Environment]::GetFolderPath("Desktop")
# remove the shortcut from the desktop if it exists
# Check if the shortcut exists
if (Test-Path $desktop_dest\manually-start-wsl-backup.lnk) {
    # If the shortcut exists, delete it
    Remove-Item -Path $desktop_dest\manually-start-wsl-backup.lnk -Force
}

# Copy the shortcut to the user's desktop

Copy-Item -Path $ShortcutPath -Destination $desktop_dest







Add-Type -AssemblyName PresentationFramework
$msgBoxInput = [System.Windows.MessageBox]::Show('The WSL backup solution has been installed' + [Environment]::NewLine + 'Click Yes to run the backup now' + [Environment]::NewLine + 'Click No to close', 'Confirmation', 'YesNo', 'Information')

switch ($msgBoxInput) {
'Yes' {
    # Run the task
    Start-ScheduledTask -TaskName $TaskName
    Start-Sleep -Seconds 10
    Get-TaskStatus -TaskName "$TaskName"
}
'No' {
    # Do nothing
}
}

# Stop the transcript  
Stop-Transcript
exit

}
catch {
    Write-Host "Error occurred: $_"
    # Handle the error (e.g., log it, display a message, etc.)
    Stop-Transcript
}



