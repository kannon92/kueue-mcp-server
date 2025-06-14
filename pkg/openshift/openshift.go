package openshift

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
	cmd := exec.CommandContext(ctx, "oc", args...)
	return cmd.CombinedOutput()
}

// MustGather runs `oc adm must-gather` with optional destination directory and
// additional arguments.
func MustGather(ctx context.Context, destDir string, extra []string) (string, error) {
	args := []string{"adm", "must-gather", "--image=registry.redhat.io/kueue-tech-preview/kueue-must-gather-rhel9:eb54d00af3f19663d025db9887d9425988cfef19"}
	if destDir != "" {
		args = append(args, fmt.Sprintf("--dest-dir=%s", destDir))
	}
	args = append(args, extra...)
	out, err := Run(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("oc adm must-gather failed: %w: %s", err, out)
	}
	return string(out), nil
}
