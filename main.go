package main

import (
	"context"
	"fmt"

	"github.com/kannon92/kueue-mcp-server/pkg/kueue"
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

// explainJobScheduling defines the explain_job_scheduling MCP tool.
var explainJobScheduling = mcp.NewTool(
	"explain_job_scheduling",
	mcp.WithTitleAnnotation("Explain why my job is not scheduled"),
	mcp.WithDescription(`Gathers information about why a job is not scheduled.`),
	mcp.WithString("namespace",
		mcp.Description("Directory to write gathered data"),
	),
	mcp.WithString("job_name",
		mcp.Description("Name of the job to explain scheduling for"),
	),
)

// explainJobScheduling defines the explain_job_scheduling MCP tool.
var explainLocalQueue = mcp.NewTool(
	"local_queue_status",
	mcp.WithTitleAnnotation("Why no workloads are admitted in the LocalQueue"),
	mcp.WithDescription(`Provides information about why no workloads are admitted in the LocalQueue.`),
	mcp.WithString("namespace",
		mcp.Description("Directory to write gathered data"),
	),
	mcp.WithString("local_queue_name",
		mcp.Description("Name of the LocalQueue to explain scheduling for"),
	),
)

// explainJobScheduling defines the explain_job_scheduling MCP tool.
var explainClusterQueue = mcp.NewTool(
	"cluster_queue_status",
	mcp.WithTitleAnnotation("Why no workloads are admitted in the ClusterQueue"),
	mcp.WithDescription(`Provides information about why no workloads are admitted in the ClusterQueue.`),
	mcp.WithString("cluster_queue_name",
		mcp.Description("Name of the ClusterQueue"),
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

// handleMustGather executes oc adm must-gather with the provided arguments.
func handleClusterQueue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cqName := req.GetString("cluster_queue_name", "")
	out, err := kueue.DescribeClusterQueue(ctx, cqName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(out), nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Kueue Support",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	RegisterTools(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// RegisterTools registers all available tools with the provided server.
func RegisterTools(s *server.MCPServer) {
	s.AddTools(
		server.ServerTool{Tool: mustGatherTool, Handler: handleMustGather},
		server.ServerTool{Tool: explainClusterQueue, Handler: handleClusterQueue},
	)
}
