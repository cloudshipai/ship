# VHS documentation tape for ship terraform-tools generate-docs
# https://github.com/charmbracelet/vhs

Output terraform-generate-docs.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 600
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Generate documentation for Terraform modules"
Enter
Sleep 1s

Type "ship terraform-tools generate-docs examples/terraform/easy-s3-bucket"
Enter
Sleep 5s

Type "# View the generated documentation"
Enter
Sleep 500ms

Type "cat examples/terraform/easy-s3-bucket/README.md | head -20"
Enter
Sleep 3s

Type "# Documentation generated successfully! 📚"
Enter
Sleep 2s