# VHS documentation tape for OpenInfraQuote demo
# https://github.com/charmbracelet/vhs

Output openinfraquote-demo.gif

# Configure terminal
Set FontSize 13
Set Width 1200
Set Height 700
Set Theme "catppuccin-mocha"
Set TypingSpeed 50ms

# Start recording
Type "# Ship CLI - OpenInfraQuote Cost Analysis 💰"
Enter
Sleep 1s

Type "# View the infrastructure we'll analyze"
Enter
Type "cat examples/terraform/simple-demo/main.tf | grep -A5 'resource \"aws'"
Enter
Sleep 3s

Type "# OpenInfraQuote analyzes Terraform plans for accurate cost estimation"
Enter
Type "# It provides detailed AWS pricing based on resource configurations"
Enter
Sleep 2s

Type "# Example usage (requires terraform plan):"
Enter
Type "# ship terraform-tools cost-analysis examples/terraform/simple-demo"
Enter
Sleep 2s

Type "# OpenInfraQuote features:"
Enter
Type "# ✓ Accurate AWS pricing data"
Enter
Type "# ✓ Supports 100+ AWS resource types" 
Enter
Type "# ✓ Region-specific pricing"
Enter
Type "# ✓ JSON output for automation"
Enter
Sleep 3s

Type "# For a working example, OpenInfraQuote would show:"
Enter
Type "# - 3x t3.medium instances: ~$91.98/month"
Enter
Type "# - 1x RDS db.t3.medium: ~$71.52/month"
Enter
Type "# - S3 storage costs based on usage"
Enter
Type "# - ALB costs: ~$16.20/month base"
Enter
Sleep 3s