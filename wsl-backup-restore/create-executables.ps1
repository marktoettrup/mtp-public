
# Install ps2exe in ps desktop
# can't run in shell: pwsh (Means powershell core)
# runs only in shell: powershell
# https://stackoverflow.com/questions/75008756/running-ps2exe-in-github-actions
# from https://www.powershellgallery.com/packages/ps2exe/1.0.10
# Start-Process pwsh -Verb RunAs
# start-process powershell -verb runas
# $PSVersionTable.PSVersion
# $PSVersionTable.PSEdition

Install-Module ps2exe  -Force
Set-Location $PSScriptRoot
Invoke-ps2exe .\setup.ps1 .\setup.exe -icon .\tux.ico
Invoke-ps2exe .\manually-start-wsl-backup.ps1 .\manually-start-wsl-backup.exe -icon .\tux.ico
Invoke-ps2exe .\manually-restore-wsl.ps1 .\manually-restore-wsl.exe -icon .\tux.ico
Invoke-ps2exe .\uninstall.ps1 .\uninstall.exe -icon .\tux.ico