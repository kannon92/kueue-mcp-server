package mcp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kannon92/kueue-mcp-server/pkg/openshift"
	mcp "github.com/mark3labs/mcp-go/mcp"
)

// mockCallToolRequest implements mcp.CallToolRequest for testing.
type mockCallToolRequest struct {
	args map[string]any
}

func (m *mockCallToolRequest) GetString(key, def string) string {
	val, ok := m.args[key]
	if !ok {
		return def
	}
	if s, ok := val.(string); ok {
		return s
	}
	return def
}

func (m *mockCallToolRequest) GetArguments() map[string]any {
	return m.args
}

func Test_handleMustGather_success(t *testing.T) {
	origRun := openshift.Run
	defer func() { openshift.Run = origRun }()
	openshift.Run = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte("must-gather output"), nil
	}

	req := &mockCallToolRequest{
		args: map[string]any{
			"dest_dir":   "/tmp/test",
			"extra_args": []any{"--foo", "--bar"},
		},
	}

	result, err := handleMustGather(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.Type != mcp.ToolResultText {
		t.Fatalf("expected ToolResultText, got %v", result)
	}
	if result.Text != "must-gather output" {
		t.Errorf("unexpected output: %q", result.Text)
	}
}

func Test_handleMustGather_error(t *testing.T) {
	origRun := openshift.Run
	defer func() { openshift.Run = origRun }()
	openshift.Run = func(ctx context.Context, args ...string) ([]byte, error) {
		return nil, errors.New("fail")
	}

	req := &mockCallToolRequest{
		args: map[string]any{
			"dest_dir":   "",
			"extra_args": []any{},
		},
	}

	result, err := handleMustGather(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.Type != mcp.ToolResultError {
		t.Fatalf("expected ToolResultError, got %v", result)
	}
	if result.Error == "" {
		t.Error("expected error message in result")
	}
}

func Test_handleMustGather_extraArgsNil(t *testing.T) {
	origRun := openshift.Run
	defer func() { openshift.Run = origRun }()
	openshift.Run = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}

	req := &mockCallToolRequest{
		args: map[string]any{
			"dest_dir": "/tmp/test",
			// "extra_args" missing
		},
	}

	result, err := handleMustGather(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.Type != mcp.ToolResultText {
		t.Fatalf("expected ToolResultText, got %v", result)
	}
	if result.Text != "ok" {
		t.Errorf("unexpected output: %q", result.Text)
	}
}

func Test_handleMustGather_extraArgsTypeMismatch(t *testing.T) {
	origRun := openshift.Run
	defer func() { openshift.Run = origRun }()
	openshift.Run = func(ctx context.Context, args ...string) ([]byte, error) {
		// Should still be called with empty extras
		if !reflect.DeepEqual(args, []string{"adm", "must-gather", "--image=registry.redhat.io/kueue-tech-preview/kueue-must-gather-rhel9:eb54d00af3f19663d025db9887d9425988cfef19", "--dest-dir=/tmp/test"}) {
			t.Errorf("unexpected args: %v", args)
		}
		return []byte("ok"), nil
	}

	req := &mockCallToolRequest{
		args: map[string]any{
			"dest_dir":   "/tmp/test",
			"extra_args": "not-an-array",
		},
	}

	result, err := handleMustGather(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.Type != mcp.ToolResultText {
		t.Fatalf("expected ToolResultText, got %v", result)
	}
	if result.Text != "ok" {
		t.Errorf("unexpected output: %q", result.Text)
	}
}
