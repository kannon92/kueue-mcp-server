# Kueue MCP Server

## Overview

**Kueue MCP Server** is a Go application that provides a set of [Mark3 Labs MCP](https://github.com/mark3labs/mcp-go) tools for OpenShift and Kueue troubleshooting and support automation. It exposes tools via an MCP server, allowing automated collection of cluster diagnostics and explanations for job and queue scheduling issues.

## Features

- **Must-Gather Tool**: Runs `oc adm must-gather` with custom arguments to collect comprehensive cluster diagnostics.
- **Explain Job Scheduling**: (Stub) Intended to explain why a specific job is not scheduled.
- **Local Queue Status**: (Stub) Intended to provide information about workloads in a LocalQueue.
- **Cluster Queue Status**: Explains why no workloads are admitted in a ClusterQueue by running `kubectl describe clusterqueue`.

## Usage

### Building

To build the server binary:

```sh
make
# or manually:
go build -o kueue-mcp-server main.go
```

### Running

Start the server (by default, it uses stdio for communication):

```sh
./kueue-mcp-server
```

The server will listen for MCP tool requests via stdio.

## Code Structure

- **main.go**: Entry point. Registers tools and starts the MCP server.
- **pkg/openshift/openshift.go**: Implements the must-gather logic using `oc`.
- **pkg/kueue/kueue.go**: Implements cluster queue description using `kubectl`.
- **pkg/mcp/mcp.go**: (If present) May contain additional MCP tool definitions and handlers.

## Extending

- Add new tools by defining them with `mcp.NewTool` and registering them in `RegisterTools`.
- Implement handlers to execute logic and return results via `mcp.CallToolResult`.

## Requirements

- Go 1.20+
- Access to `oc` and `kubectl` CLIs in the environment where the server runs.
- OpenShift/Kubernetes cluster access (for must-gather and describe commands).

## Example Tools

### Must-Gather

Collects cluster data for support:

- **Parameters**:
  - `dest_dir`: Directory to write gathered data.
  - `extra_args`: Additional arguments for `oc adm must-gather`.

### Cluster Queue Status

Explains why no workloads are admitted in a ClusterQueue:

- **Parameters**:
  - `cluster_queue_name`: Name of the ClusterQueue.

## Integration With Claude Desktop

To use the Kueue MCP Server as a custom MCP tool in your Claude Desktop configuration:

1. **Build and run the server**  
   Make sure the `kueue-mcp-server` binary is built and available on your system. Start the server so it listens for stdio requests:

   ```sh
   ./kueue-mcp-server
   ```

2. **Add the tool to your Claude Desktop config**  
   - For Fedora, the following yaml could be used:
    I found the config in ~/.config/Claude/claude_desktop_config.json

```yaml
   {
        "mcpServers": {
                "kubernetes": {
                        "command": "npx",
                        "args": [
                                "-y",
                                "kubernetes-mcp-server@latest"
                        ]
        },
        "kueue-mcp-server": {
            "command": "/path/to/kueue-mcp-server"
  }
  
}
```
   - Replace `/path/to/kueue-mcp-server` with the full path to your built binary.
   - Ensure `stdin` and `stdout` are enabled so Claude Desktop can communicate with the server.

3. **Restart Claude Desktop**  
   After updating the configuration, restart Claude Desktop to load the new tool.

4. **Use the tool**  
   The Kueue MCP Server tools (such as must-gather and cluster queue status) will now be available in Claude Desktop’s tool palette.
