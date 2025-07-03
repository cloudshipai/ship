# VHS documentation tape for ship terraform-tools checkov-scan
# https://github.com/charmbracelet/vhs

Output terraform-checkov-scan.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 800
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Comprehensive security scan with Checkov"
Enter
Sleep 1s

Type "ship terraform-tools checkov-scan /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket"
Enter
Sleep 7s

Type "# Checkov found several security checks - 8 passed, 5 failed"
Enter
Sleep 2s