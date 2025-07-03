# VHS documentation tape for OpenInfraQuote real demo
# https://github.com/charmbracelet/vhs

Output openinfraquote-real-demo.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 800
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Ship CLI - OpenInfraQuote Cost Analysis 💰"
Enter
Sleep 1s

Type "# Analyze a Terraform plan with OpenInfraQuote"
Enter
Type "ship terraform-tools cost-analysis test/terraform-mocks/tfplan.json"
Enter
Sleep 6s

Type "# OpenInfraQuote provides detailed cost breakdown:"
Enter
Type "# ✓ EC2 t3.medium instance: $30.37/month ($0.0416/hour)"
Enter
Type "# ✓ S3 bucket operations: ~$40/month (estimated)"
Enter
Type "# ✓ Total estimated cost: $70.37/month"
Enter
Sleep 3s

Type "# Key features:"
Enter
Type "# - Real AWS pricing data updated daily"
Enter
Type "# - Accounts for all resource attributes"
Enter
Type "# - JSON output for CI/CD integration"
Enter
Sleep 2s