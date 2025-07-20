# Uninstall Script

# Function to remove PowerShell 7
# https://github.com/PowerShell/PowerShell/releases/

# Function to remove scheduled tasks
function Remove-ScheduledTasks {
    Write-Host "Removing scheduled tasks..."
    $taskNames = @("Backup WSL", "Backup WSL Silent")
    foreach ($taskName in $taskNames) {
        if (Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue) {
            Unregister-ScheduledTask -TaskName $taskName -Confirm:$false
            Write-Host "Scheduled task '$taskName' removed."
        } else {
            Write-Host "Scheduled task '$taskName' not found."
        }
    }
}

# Function to remove backup files
function Remove-BackupScriptFiles {
    Write-Host "Removing backup script files..."
    $backupPath = "$env:OneDrive\wsl-scripts"
    if (Test-Path $backupPath) {
        Remove-Item -Recurse -Force $backupPath
        Write-Host "Backup files removed from $backupPath."
    } else {
        Write-Host "No backup files found in $backupPath."
    }
}



# Remove scheduled tasks
Remove-ScheduledTasks

# Remove backup files
Remove-BackupScriptFiles

# Clean up any other resources
Write-Host "Cleaning up other resources..."
# Add any other cleanup logic here

Write-Host "Uninstallation complete."
