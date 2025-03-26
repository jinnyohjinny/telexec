package controller

import (
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

type CmdExecResult struct {
	Stdout []byte
	Stderr []byte
}

func NewCmdOutputWriter(timeout int, workDir string) *CmdOutputWriter {
	return &CmdOutputWriter{
		TimeoutSecond: timeout,
		WorkDir:       workDir,
	}
}

func (c *CmdOutputWriter) ExecOutput(command string) (outOk, outErr []byte, err error) {
	chResp := make(chan CmdExecResult)

	executor := "/bin/bash"
	cmd := exec.Command(executor, "-c", command)

	if c.WorkDir != "" {
		cmd.Dir = c.WorkDir
	}

	go func(chn chan CmdExecResult, cmd *exec.Cmd) {
		c.mu.Lock()
		defer c.mu.Unlock()
		out, err := cmd.CombinedOutput()
		errOut := []byte("")
		if err != nil {
			errOut = []byte(err.Error())
		}

		chn <- CmdExecResult{Stdout: out, Stderr: errOut}
	}(chResp, cmd)

	if c.TimeoutSecond == 0 {
		result := <-chResp
		outOk = result.Stdout
		outErr = result.Stderr
		return
	}

	select {
	case result := <-chResp:
		outOk = result.Stdout
		outErr = result.Stderr
		cmd.Process.Kill()
	case <-time.After(time.Second * time.Duration(c.TimeoutSecond)):
		cmd.Process.Kill()
		errMsg := "TIMEOUT: EXCEEDED -- Process killed"
		err = errors.New(errMsg)
		outErr = []byte(errMsg)
	}
	return
}

func (c *CmdOutputWriter) ExecHeadOutput(command string) (outOk, outErr []byte, err error) {
	return c.ExecOutput(fmt.Sprintf("%s | head", command))
}

func (c *CmdOutputWriter) ExecTailOutput(command string) (outOk, outErr []byte, err error) {
	return c.ExecOutput(fmt.Sprintf("%s | tail", command))
}
