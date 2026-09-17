package runner

import (
	"testing"
)

func TestInterpolateCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "No args",
			command:  "echo hello",
			args:     nil,
			expected: "echo hello",
		},
		{
			name:     "Positional $1 and $2",
			command:  "docker build -t app:$1 $2",
			args:     []string{"v1.0", "./src"},
			expected: "docker build -t app:v1.0 ./src",
		},
		{
			name:     "Braced placeholders {{1}}",
			command:  "git commit -m '{{1}}'",
			args:     []string{"initial commit"},
			expected: "git commit -m 'initial commit'",
		},
		{
			name:     "All args $@",
			command:  "npm run test -- $@",
			args:     []string{"--watch", "--verbose"},
			expected: "npm run test -- --watch --verbose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InterpolateCommand(tt.command, tt.args)
			if got != tt.expected {
				t.Errorf("InterpolateCommand(%q, %v) = %q, want %q", tt.command, tt.args, got, tt.expected)
			}
		})
	}
}

func TestExecuteCommands_DryRunWithInterpolation(t *testing.T) {
	commands := []string{"echo $1", "echo $2"}
	opts := ExecuteOptions{
		Args:   []string{"argA", "argB"},
		DryRun: true,
	}

	err := ExecuteCommands(commands, opts)
	if err != nil {
		t.Fatalf("expected dry run to succeed, got %v", err)
	}
}

func TestExecuteCommands_SequentialSuccess(t *testing.T) {
	commands := []string{"echo runner1", "echo runner2"}
	opts := ExecuteOptions{
		SkipPrompt: true,
	}

	err := ExecuteCommands(commands, opts)
	if err != nil {
		t.Fatalf("expected sequential execution to succeed, got %v", err)
	}
}

func TestExecuteCommands_SequentialFailFast(t *testing.T) {
	commands := []string{"exit 1", "echo should-not-run"}
	opts := ExecuteOptions{
		SkipPrompt:      true,
		ContinueOnError: false,
	}

	err := ExecuteCommands(commands, opts)
	if err == nil {
		t.Fatal("expected execution to fail when command exits non-zero")
	}
}

func TestExecuteCommands_ParallelSuccess(t *testing.T) {
	commands := []string{"echo p1", "echo p2"}
	opts := ExecuteOptions{
		Parallel: true,
	}

	err := ExecuteCommands(commands, opts)
	if err != nil {
		t.Fatalf("expected parallel execution to succeed, got %v", err)
	}
}

