package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type CmdOutputWriter struct {
	TimeoutSecond int
	WorkDir       string
	mu            sync.Mutex
}

func NewCmdOutputWriter(timeout int, workDir string) *CmdOutputWriter {
	return &CmdOutputWriter{
		TimeoutSecond: timeout,
		WorkDir:       workDir,
	}
}

func (c *CmdOutputWriter) ExecOutput(command string) ([]byte, []byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeoutSecond)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	if c.WorkDir != "" {
		cmd.Dir = c.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, []byte("ERROR: Command timed out"), ctx.Err()
	}

	if stdout.Len() == 0 && stderr.Len() == 0 && err == nil {
		stdout.WriteString("Command executed successfully (no output)")
	}

	return stdout.Bytes(), stderr.Bytes(), err
}

func (c *CmdOutputWriter) ExecHeadOutput(command string) ([]byte, []byte, error) {
	return c.ExecOutput(fmt.Sprintf("%s | head -n 20", command))
}

func (c *CmdOutputWriter) ExecTailOutput(command string) ([]byte, []byte, error) {
	return c.ExecOutput(fmt.Sprintf("%s | tail -n 20", command))
}
