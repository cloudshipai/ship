# VHS documentation tape for ship terraform-tools security-scan
# https://github.com/charmbracelet/vhs

Output terraform-security-scan.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 800
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Security scan Terraform code with Trivy"
Enter
Sleep 1s

Type "ship terraform-tools security-scan /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket"
Enter
Sleep 6s

Type "# Security scan found 2 issues to review"
Enter
Sleep 2s