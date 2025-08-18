const { spawn } = require('child_process');

// Test the MCP server with a direct tool call
const mcp = spawn('./ship', ['mcp', 'custodian']);

// Send initialization
const initRequest = {
  jsonrpc: "2.0",
  method: "initialize",
  params: {
    protocolVersion: "2024-11-05",
    capabilities: {},
    clientInfo: { name: "test-client", version: "1.0.0" }
  },
  id: 1
};

mcp.stdin.write(JSON.stringify(initRequest) + '\n');

// Wait for initialization response, then send tool call
setTimeout(() => {
  const toolCall = {
    jsonrpc: "2.0",
    method: "tools/call",
    params: {
      name: "custodian_get_version",
      arguments: {}
    },
    id: 2
  };
  
  console.log('Sending tool call:', JSON.stringify(toolCall));
  mcp.stdin.write(JSON.stringify(toolCall) + '\n');
}, 1000);

// Handle output
mcp.stdout.on('data', (data) => {
  const lines = data.toString().split('\n').filter(line => line.trim());
  lines.forEach(line => {
    try {
      const response = JSON.parse(line);
      if (response.id === 2) {
        console.log('Tool response:', JSON.stringify(response, null, 2));
        process.exit(0);
      }
    } catch (e) {
      // Ignore non-JSON lines
    }
  });
});

mcp.stderr.on('data', (data) => {
  console.error('MCP stderr:', data.toString());
});

// Timeout after 30 seconds
setTimeout(() => {
  console.log('Timeout - no response received');
  mcp.kill();
  process.exit(1);
}, 30000);