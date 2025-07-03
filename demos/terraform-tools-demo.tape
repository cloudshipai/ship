# VHS documentation tape for all ship terraform-tools
# https://github.com/charmbracelet/vhs

Output terraform-tools-demo.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 800
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Ship CLI Terraform Tools Demo 🚢"
Enter
Sleep 1s

Type "# Working with example S3 bucket module"
Enter
Type "ls /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket/"
Enter
Sleep 2s

Type "# 1. Generate documentation"
Enter
Type "ship terraform-tools generate-docs /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket"
Enter
Sleep 4s

Type "# 2. Lint the code"
Enter
Type "ship terraform-tools lint /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket"
Enter
Sleep 4s

Type "# 3. Security scan with Trivy"
Enter
Type "ship terraform-tools security-scan /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket | head -30"
Enter
Sleep 5s

Type "# 4. Comprehensive scan with Checkov"
Enter
Type "ship terraform-tools checkov-scan /home/epuerta/projects/ship/examples/terraform/easy-s3-bucket | grep -A5 'results':"
Enter
Sleep 4s

Type "# All tools executed successfully! 🎉"
Enter
Sleep 2s