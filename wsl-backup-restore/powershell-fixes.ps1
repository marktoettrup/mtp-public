# Fix powershell right click context menu
# run elevated
# from https://github.com/PowerShell/PowerShell/issues/14216#issuecomment-1820020030

# author: eterna1_0blivion & soredake

# Get a link to the Registration Entries file
$link = "https://gist.github.com/eterna1-0blivion/70c1e5b14c7cfa8c6b6d574eb38fd27e/raw/context_pwsh_fix.reg"

# Install them into Registry using PowerShell
Invoke-WebRequest -Uri "$link" -OutFile "$env:TEMP/context_pwsh_fix.reg"
reg import "$env:TEMP/context_pwsh_fix.reg"

# Notify the user before exiting the program
Read-Host -Prompt "Press Enter to exit"
Exit
