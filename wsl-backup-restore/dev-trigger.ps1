$BackupDir = "$env:OneDrive\wsl-backup"
$FilePath = "$env:OneDrive\wsl-scripts"
$ScriptPath = "$FilePath\backup-wsl.ps1"
# pwsh.exe -ExecutionPolicy Bypass -File $ScriptPath silentbackup $BackupDir

# to trigger the silent backup
. $ScriptPath -Command 'silentbackup' -BackupDirectory $BackupDir

# to trigger the normal backup with user interaction
. $ScriptPath -Command 'backup' -BackupDirectory $BackupDir
