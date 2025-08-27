FROM mcp/grafana:latest

# Set working directory
WORKDIR /workspace

# Default command
CMD ["mcp-grafana", "--help"]