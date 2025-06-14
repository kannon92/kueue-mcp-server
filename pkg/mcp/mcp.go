package mcp

import (
	"context"
	"fmt"

	"github.com/kannon92/kueue-mcp-server/pkg/openshift"
	mcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// mustGatherTool defines the collect_must_gather MCP tool.
var mustGatherTool = mcp.NewTool(
	"collect_must_gather_ocp_kueue",
	mcp.WithTitleAnnotation("Collect cluster data via oc adm must-gather"),
	mcp.WithDescription(`Runs "oc adm must-gather" to capture debugging information.

Create a temporary directory and pass it using the --dest-dir option to store the output in a single place.

oc adm must-gather can scoop up almost every artifact engineers or support need in a single shot: it exports the full YAML for all cluster-scoped and namespaced resources (Deployments, CRDs, Nodes, ClusterOperators, etc.); captures pod and container logs as well as systemd journal slices from each node to trace runtime crashes or OOMs; grabs API-server and OAuth audit logs for security or compliance forensics; collects kernel, cgroup, and other node sysinfo plus tuned and kubelet configs for performance tuning; optionally runs add-on scripts such as gather_network_logs to archive iptables/OVN flows and CNI pod logs, or gather_profiling_node to fetch 30-second CPU and heap pprof dumps from both kubelet and CRI-O for hotspot analysis; and, through plug-in images, can extend to operator-specific data like storage states or virtualization metrics, ensuring one reproducible tarball contains configuration, logs, network traces, performance profiles, and security audits for thorough offline debugging. Use "oc adm must-gather -h" for available options.`),
	mcp.WithString("dest_dir",
		mcp.Description("Directory to write gathered data"),
	),
	mcp.WithArray("extra_args",
		mcp.Description("Additional arguments passed directly to oc adm must-gather"),
		mcp.Items(map[string]any{"type": "string"}),
	),
)

// handleMustGather executes oc adm must-gather with the provided arguments.
func handleMustGather(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dest := req.GetString("dest_dir", "")
	extraAny, _ := req.GetArguments()["extra_args"].([]any)
	extras := make([]string, len(extraAny))
	for i, a := range extraAny {
		extras[i] = fmt.Sprint(a)
	}
	out, err := openshift.MustGather(ctx, dest, extras)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(out), nil
}

// RegisterTools registers all available tools with the provided server.
func RegisterTools(s *server.MCPServer) {
	s.AddTools(
		server.ServerTool{Tool: mustGatherTool, Handler: handleMustGather},
	)
}
