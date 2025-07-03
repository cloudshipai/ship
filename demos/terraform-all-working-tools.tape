# VHS documentation tape for working ship terraform-tools
# https://github.com/charmbracelet/vhs

Output terraform-all-working-tools.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 700
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Ship CLI Terraform Tools Suite 🚢"
Enter
Sleep 1s

Type "# Generate documentation for Terraform module"
Enter
Type "ship terraform-tools generate-docs examples/terraform/easy-s3-bucket"
Enter
Sleep 4s

Type "# Lint Terraform code"
Enter
Type "ship terraform-tools lint examples/terraform/easy-s3-bucket"
Enter
Sleep 4s

Type "# Security scan with Trivy"
Enter
Type "ship terraform-tools security-scan examples/terraform/easy-s3-bucket --push"
Enter
Sleep 6s

Type "# All tools executed successfully! 🎉"
Enter
Sleep 2s