package kueue

import (
	"context"
	"fmt"
	"os/exec"
)

// run executes the oc command with given arguments and returns combined output.
//
// It is defined as a variable to allow tests to substitute a fake implementation
// without spawning external processes.
// Run is used by helper functions to execute the oc command. Tests may override
// this variable to avoid running external commands.
var Run = func(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	return cmd.CombinedOutput()
}

// DescribeClusterQueue describes a ClusterQueue using `kubectl describe clusterqueue <name>`.
func DescribeClusterQueue(ctx context.Context, clusterQueue string) (string, error) {
	args := []string{"describe", "clusterqueue", clusterQueue}
	out, err := Run(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("kubectl describe cluster queue failed: %w: %s", err, out)
	}
	return string(out), nil
}
