<#
.SYNOPSIS
    Restores a WSL distribution from a backup file.

.DESCRIPTION
    This script guides the user through restoring a WSL distribution from a backup TAR file.
    The user is prompted to select the backup file, provide a new distribution name,
    and choose a destination folder. It then uses the WSL import command to restore the distribution.
    Notifications are provided via the BurntToast module.

.NOTES
    Tested with PowerShell 7.5.0. Ensure that PowerShell 7.4.1 or later is installed.
    For more info: https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows

.LINK
    https://lazyadmin.nl/powershell/how-to-create-a-powershell-scheduled-task/
#>

Set-Culture en-US

# Set log file
$logfile = Join-Path "$env:OneDrive\wsl-scripts" "$(Hostname)-manually-restore-wsl.log"
Start-Transcript -Path $logfile

try {
    # Load BurntToast if available
    try {
        Import-Module -Name BurntToast -ErrorAction Stop
    } catch {
        Write-Host "BurntToast module not available. Skipping toast notifications."
    }

    # Load assemblies for GUI dialogs
    Add-Type -AssemblyName System.Windows.Forms
    Add-Type -AssemblyName Microsoft.VisualBasic

    # Set backup folder path
    $backupFolder = Join-Path $env:OneDrive "wsl-backup"
    if (-not (Test-Path $backupFolder)) {
        throw "Backup folder not found: $backupFolder"
    }

    # Prompt user for backup file
    $openFileDialog = New-Object System.Windows.Forms.OpenFileDialog
    $openFileDialog.InitialDirectory = $backupFolder
    $openFileDialog.Filter = "Backup Files (*.tar)|*.tar|All Files (*.*)|*.*"
    $openFileDialog.Title = "Select the WSL Backup File to Restore"

    if ($openFileDialog.ShowDialog() -ne [System.Windows.Forms.DialogResult]::OK) {
        throw "User cancelled the file selection."
    }

    $backupFile = $openFileDialog.FileName
    Write-Host "Selected backup file: $backupFile"

    # Prompt for distro name
    $distroName = [Microsoft.VisualBasic.Interaction]::InputBox(
        "Enter the new WSL distribution name:", "Distribution Name", "RestoredWSL"
    ).Trim()
    if ([string]::IsNullOrWhiteSpace($distroName)) {
        throw "No distribution name provided."
    }

    # Prompt for destination folder
    $defaultDestination = "E:\\RestoredWSL"
    $destinationFolder = [Microsoft.VisualBasic.Interaction]::InputBox(
        "Enter the destination folder for the restored distribution:",
        "Destination Folder",
        $defaultDestination
    ).Trim()
    if ([string]::IsNullOrWhiteSpace($destinationFolder)) {
        throw "No destination folder provided."
    }

    Write-Host "Destination folder: $destinationFolder"

    # Confirm restore operation
    $confirmation = [System.Windows.Forms.MessageBox]::Show(
        "Proceed with restoring WSL distribution '$distroName' from backup file:`n$backupFile`nto destination:`n$destinationFolder",
        "Confirm Restore",
        [System.Windows.Forms.MessageBoxButtons]::YesNo,
        [System.Windows.Forms.MessageBoxIcon]::Question
    )
    if ($confirmation -ne [System.Windows.Forms.DialogResult]::Yes) {
        throw "User cancelled the restore operation."
    }

    # Create or confirm overwrite of destination folder
    if (Test-Path $destinationFolder) {
        $overwriteConfirmation = [System.Windows.Forms.MessageBox]::Show(
            "Destination folder already exists. Overwrite it?",
            "Confirm Overwrite",
            [System.Windows.Forms.MessageBoxButtons]::YesNo,
            [System.Windows.Forms.MessageBoxIcon]::Warning
        )
        if ($overwriteConfirmation -ne [System.Windows.Forms.DialogResult]::Yes) {
            throw "User cancelled due to existing destination folder."
        }
        Remove-Item -Path $destinationFolder -Recurse -Force
    } else {
        New-Item -ItemType Directory -Path $destinationFolder | Out-Null
    }

    # Execute WSL import
    Write-Host "Restoring WSL distribution..."
    $originalEncoding = [console]::OutputEncoding
    [console]::OutputEncoding = New-Object System.Text.UnicodeEncoding
    wsl.exe --import "$distroName" "$destinationFolder" "$backupFile"
    [console]::OutputEncoding = $originalEncoding

    # Notify success
    New-BurntToastNotification -Text "WSL Restore", "Distribution '$distroName' restored successfully."
    Write-Host "WSL distribution '$distroName' restored successfully."

} catch {
    Write-Host "Error occurred: $_"
    New-BurntToastNotification -Text "WSL Restore Failed", "An error occurred: $_"
} finally {
    Stop-Transcript
}
