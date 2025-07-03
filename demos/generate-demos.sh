#!/bin/bash

# Generate all Ship CLI demo GIFs

echo "🎬 Generating Ship CLI demo GIFs..."
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# List of all tape files
TAPES=(
    "ship-quick-start.tape"
    "ship-auth.tape"
    "ship-query.tape"
    "ship-ai-investigate.tape"
    "ship-ai-agent.tape"
    "ship-push.tape"
    "ship-mcp.tape"
    "ship-terraform-all-tools.tape"
)

# Generate each GIF
for tape in "${TAPES[@]}"; do
    if [ -f "$tape" ]; then
        echo "📹 Recording $tape..."
        vhs "$tape"
        echo "✅ Generated ${tape%.tape}.gif"
        echo ""
    else
        echo "❌ Warning: $tape not found"
    fi
done

echo "🎉 All demos generated!"
echo ""
echo "📁 Files created:"
ls -la *.gif 2>/dev/null | awk '{print "   - " $9}'
echo ""
echo "📖 View the demo gallery: README.md"