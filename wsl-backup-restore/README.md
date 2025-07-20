# WSL Backup & Restore Toolkit

A user-friendly PowerShell-based solution to **backup, restore, and schedule Windows Subsystem for Linux (WSL)** distributions. This toolkit provides manual and scheduled options for creating `.tar` or `.7z`archives of your WSL environments and restoring them when needed.

## 🧰 Features

- ✅ Will backup all of the WSLs you have installed
- ✅ Options: 
  - ✅ scheduled backup of all WSLs (as a user you will be prompted every 14th day to do the backup, just click OK)
  - ✅ You can trigger manual backup of all WSLs at any time
  - ✅ you will get a scheduled silent backup of all WSLs that runs every Sunday at 10:00 AM
  - ✅ Change the schedules of the backups to suit your needs
- ✅ The solution saves your WSL backups as 7z files if you have 7z on your Windows host, otherwise .tar file backups are saved, but they are often three times larger, so be careful. The Backups are saved to your One-Drive to ensure online backup (great way to use One-drive for something useful)
- ✅ Rotates your backups, keeping the last five (5) backups in One-drive
- ✅ Simple/easy restore from One-drive
- ✅ Just install and don't worry

---

## 🛠️ Installation

1. **Clone the repository** to your local machine.
2. **Activate _Admin by Request_** (Systematic's elevation software).
3. Open an **elevated PowerShell terminal**.
4. Run:
   ```powershell
   ./setup.ps1
   ```

This will install required modules and set up two scheduled tasks:

- `Backup WSL`: Runs every second Monday at 10:00 AM, with full user interaction.
- `Backup WSL Silent`: Runs silently every Sunday at 10:00 AM without user prompts.

---

## 📦 Scripts Overview

| Script Name                         | Description |
|------------------------------------|-------------|
| `backup-wsl.ps1`                   | Main script to perform a WSL backup |
| `backup-wsl-bootstrap.ps1`         | Bootstrap backup configuration and schedule |
| `backup-wsl-schedule-disable.ps1`  | Disables the backup schedule |
| `manually-restore-wsl.ps1`         | GUI-driven restore script using `wsl --import` |
| `manually-start-wsl-backup.ps1`    | Triggers backup manually |
| `manually-start-wsl-silentbackup.ps1` | Triggers silent backup without user interaction |
| `get-user-input-bootstrap.ps1`     | Collects input and configures backup options |
| `install-pwsh7.ps1`                | Installs PowerShell 7 (if not installed) |
| `uninstall.ps1`                    | Cleans up and removes backup setup |
| `create-executables.ps1`           | Wraps PowerShell scripts as executables |
| `powershell-fixes.ps1`             | Contains helper fixes for common PowerShell issues |
| `dev-trigger.ps1`                  | Used for development testing |
| `setup.ps1`                        | Main setup entry point |
| `tux.ico`                          | Linux penguin icon used in BurntToast notifications |

---

## 🚀 Usage

### 🔹 Backup WSL (Manual)
```powershell
./manually-start-wsl-backup.ps1
```

### 🔹 Silent Backup (Scheduled or manual)
```powershell
./manually-start-wsl-silentbackup.ps1
```

### 🔹 Restore WSL from Backup
```powershell
./manually-restore-wsl.ps1
```

> The restore script uses GUI dialogs to prompt for:
> - Backup file
> - New distribution name
> - Destination folder

### 🔹 Schedule Backup Task
You can change the backup scheduels in the backup-wsl-bootstrap.ps1, and then deploy with
```powershell
./backup-wsl-bootstrap.ps1
```

---

## 🔧 Requirements

- Windows 10/11 with WSL 2 installed
- PowerShell 7.4+ (auto-installed if not available)
- Administrator privileges for task scheduling
- [BurntToast](https://github.com/Windos/BurntToast) PowerShell module (optional)

---

## 📂 Output

Backups are stored in your OneDrive under:
```
OneDrive\wsl-backup\<MachineName>\<Timestamp>\<distro>.tar
```

---

## 🧹 Uninstall

To remove scheduled tasks and configuration:
```powershell
./uninstall.ps1
```

---

## 📢 Credits

Created by Mark Tøttrup, mtp, 2025  
Inspired by practical needs to safeguard WSL workflows and configurations with minimal friction.
