# scripts/install.ps1
$ErrorActionPreference = "Stop"
$GitHubRepo = "TheAICompanyLabs/mcp-proxy-open-source"
$AppName = "mcp-proxy"

Write-Host "📥 Installing Universal MCP Proxy for Windows..."

# Fetch latest release version
$ApiUrl = "https://api.github.com/repos/$GitHubRepo/releases/latest"
$Release = Invoke-RestMethod -Uri $ApiUrl
$Version = $Release.tag_name

if (-not $Version) {
    Write-Error "❌ Failed to fetch latest version. Ensure the repository is public."
    exit 1
}

$FileName = "$AppName-windows-amd64.zip"
$DownloadUrl = "https://github.com/$GitHubRepo/releases/download/$Version/$FileName"

Write-Host " -> Downloading $Version..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $FileName

# Create destination directory if it doesn't exist
$InstallDir = "C:\mcp-proxy"
if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host " -> Extracting to $InstallDir..."
Expand-Archive -Path $FileName -DestinationPath $InstallDir -Force

Remove-Item -Path $FileName -Force
Write-Host "✅ Installation complete!"
Write-Host "⚠️ Please add $InstallDir to your System Environment Variables (PATH) to use '$AppName' globally."
