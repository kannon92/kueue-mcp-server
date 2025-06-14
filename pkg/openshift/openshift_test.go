package openshift

import (
	"context"
	"errors"
	"testing"
)

func TestMustGather_Success(t *testing.T) {
	origRun := Run
	defer func() { Run = origRun }()

	Run = func(ctx context.Context, args ...string) ([]byte, error) {
		expected := []string{"adm", "must-gather", "--image=registry.redhat.io/kueue-tech-preview/kueue-must-gather-rhel9:eb54d00af3f19663d025db9887d9425988cfef19", "--dest-dir=/tmp/test", "--foo"}
		for i, arg := range expected {
			if args[i] != arg {
				t.Errorf("expected arg %d to be %q, got %q", i, arg, args[i])
			}
		}
		return []byte("gathered"), nil
	}

	out, err := MustGather(context.Background(), "/tmp/test", []string{"--foo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "gathered" {
		t.Errorf("expected output 'gathered', got %q", out)
	}
}

func TestMustGather_Error(t *testing.T) {
	origRun := Run
	defer func() { Run = origRun }()

	Run = func(ctx context.Context, args ...string) ([]byte, error) {
		return []byte("fail output"), errors.New("fail error")
	}

	_, err := MustGather(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got, want := err.Error(), "oc adm must-gather failed: fail error: fail output"; got != want {
		t.Errorf("unexpected error message: got %q, want %q", got, want)
	}
}
