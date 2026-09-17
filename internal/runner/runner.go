package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
)

type ExecuteOptions struct {
	Args            []string
	SkipPrompt      bool
	Parallel        bool
	ContinueOnError bool
	DryRun          bool
}

var tagColors = []string{
	"\033[36m", // Cyan
	"\033[33m", // Yellow
	"\033[32m", // Green
	"\033[35m", // Magenta
	"\033[34m", // Blue
	"\033[31m", // Red
}

const colorReset = "\033[0m"

type syncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (sw *syncWriter) PrintLine(tag, line string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Fprintf(sw.w, "%s %s\n", tag, line)
}

// ExecuteCommands runs a slice of commands sequentially or in parallel, with argument substitution.
func ExecuteCommands(commands []string, opts ExecuteOptions) error {
	if len(commands) == 0 {
		fmt.Println("No commands found in this transit.")
		return nil
	}

	// Preprocess commands with dynamic arguments
	interpolated := make([]string, len(commands))
	for i, cmd := range commands {
		interpolated[i] = InterpolateCommand(cmd, opts.Args)
	}

	if opts.DryRun {
		mode := "Sequential"
		if opts.Parallel {
			mode = "Parallel"
		}
		fmt.Printf("🔍 [Dry Run] %s execution of %d command(s):\n", mode, len(interpolated))
		for i, cmd := range interpolated {
			fmt.Printf("  %d. %s\n", i+1, cmd)
		}
		return nil
	}

	if opts.Parallel {
		return executeParallel(interpolated, opts)
	}
	return executeSequential(interpolated, opts)
}

// InterpolateCommand replaces placeholders ($1, $2, $@, {{1}}) with provided arguments.
func InterpolateCommand(command string, args []string) string {
	if len(args) == 0 {
		return command
	}

	allArgs := strings.Join(args, " ")
	res := strings.ReplaceAll(command, "$@", allArgs)
	res = strings.ReplaceAll(res, "$*", allArgs)
	res = strings.ReplaceAll(res, "{{@}}", allArgs)

	for i, arg := range args {
		pDollar := fmt.Sprintf("$%d", i+1)
		pBrace := fmt.Sprintf("{{%d}}", i+1)
		res = strings.ReplaceAll(res, pDollar, arg)
		res = strings.ReplaceAll(res, pBrace, arg)
	}
	return res
}

func executeSequential(commands []string, opts ExecuteOptions) error {
	total := len(commands)
	startTime := time.Now()

	for i, cmdStr := range commands {
		stepPrefix := fmt.Sprintf("[%d/%d]", i+1, total)

		if !opts.SkipPrompt {
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("%s Execute: %s [y/n]", stepPrefix, cmdStr),
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt cancelled: %v\n", err)
				return err
			}

			if result != "y" && result != "" {
				fmt.Printf("%s Skipped: %s\n", stepPrefix, cmdStr)
				continue
			}
		}

		fmt.Printf("\n%s ➜ Running: %s\n", stepPrefix, cmdStr)
		cmdStart := time.Now()
		err := runCommand(cmdStr)
		cmdDuration := time.Since(cmdStart).Round(time.Millisecond)

		if err != nil {
			fmt.Printf("%s ✗ Failed (%s): %v\n", stepPrefix, cmdDuration, err)
			if !opts.ContinueOnError {
				return fmt.Errorf("command %d failed: %w", i+1, err)
			}
		} else {
			fmt.Printf("%s ✓ Completed (%s)\n", stepPrefix, cmdDuration)
		}
	}

	totalDuration := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("\n✨ Finished all %d commands in %s\n", total, totalDuration)
	return nil
}

func executeParallel(commands []string, opts ExecuteOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		select {
		case sig := <-sigChan:
			fmt.Printf("\n⚠️  Received %s. Terminating all processes...\n", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	fmt.Printf("🚀 Running %d commands concurrently...\n\n", len(commands))
	startTime := time.Now()

	writer := &syncWriter{w: os.Stdout}
	var wg sync.WaitGroup
	errChan := make(chan error, len(commands))

	for idx, cmdStr := range commands {
		wg.Add(1)
		color := tagColors[idx%len(tagColors)]
		tag := fmt.Sprintf("%s[%d: %s]%s", color, idx+1, truncate(cmdStr, 20), colorReset)

		go func(cmdString string, tag string) {
			defer wg.Done()

			cmd := getShellCommand(ctx, cmdString)

			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				writer.PrintLine(tag, fmt.Sprintf("Error creating stdout pipe: %v", err))
				errChan <- err
				return
			}
			stderrPipe, err := cmd.StderrPipe()
			if err != nil {
				writer.PrintLine(tag, fmt.Sprintf("Error creating stderr pipe: %v", err))
				errChan <- err
				return
			}

			if err := cmd.Start(); err != nil {
				writer.PrintLine(tag, fmt.Sprintf("Failed to start: %v", err))
				errChan <- err
				if !opts.ContinueOnError {
					cancel()
				}
				return
			}

			var pipeWg sync.WaitGroup
			pipeWg.Add(2)

			streamOutput := func(r io.Reader) {
				defer pipeWg.Done()
				scanner := bufio.NewScanner(r)
				for scanner.Scan() {
					writer.PrintLine(tag, scanner.Text())
				}
			}

			go streamOutput(stdoutPipe)
			go streamOutput(stderrPipe)

			pipeWg.Wait()
			runErr := cmd.Wait()
			if runErr != nil {
				writer.PrintLine(tag, fmt.Sprintf("Exited with error: %v", runErr))
				errChan <- runErr
				if !opts.ContinueOnError {
					cancel()
				}
			} else {
				writer.PrintLine(tag, "Process completed successfully ✓")
			}
		}(cmdStr, tag)
	}

	wg.Wait()
	close(errChan)

	totalDuration := time.Since(startTime).Round(time.Millisecond)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		fmt.Printf("\n⚠️  Finished in %s with %d error(s)\n", totalDuration, len(errs))
		return fmt.Errorf("%d command(s) failed during parallel execution", len(errs))
	}

	fmt.Printf("\n✨ All %d commands finished successfully in %s\n", len(commands), totalDuration)
	return nil
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func getShellCommand(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		if ctx != nil {
			return exec.CommandContext(ctx, "cmd", "/C", command)
		}
		return exec.Command("cmd", "/C", command)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	if ctx != nil {
		return exec.CommandContext(ctx, shell, "-c", command)
	}
	return exec.Command(shell, "-c", command)
}

func runCommand(command string) error {
	cmd := getShellCommand(nil, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
