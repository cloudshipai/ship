# VHS documentation tape for ship terraform-tools lint
# https://github.com/charmbracelet/vhs

Output terraform-lint.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 600
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Lint Terraform code with TFLint"
Enter
Sleep 1s

Type "ship terraform-tools lint examples/terraform/easy-s3-bucket"
Enter
Sleep 5s

Type "# ✅ No linting issues found!"
Enter
Sleep 2s