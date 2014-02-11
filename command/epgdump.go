package command

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

var EpgdumpPath = "/usr/local/bin/epgdump"

type Epgdump struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func (e *Epgdump) Start() error {
	e.cmd = exec.Command(EpgdumpPath, "json", "-", "-")

	stdin, err := e.cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := e.cmd.Start(); err != nil {
		return err
	}

	e.stdin = stdin
	e.stdout = stdout
	return nil
}

func (e *Epgdump) Write(p []byte) (int, error) {
	if e.stdin == nil {
		if err := e.Start(); err != nil {
			return 0, err
		}
	}
	n, err := e.stdin.Write(p)
	return n, err
}

func (e *Epgdump) Read(p []byte) (int, error) {
	if e.stdout == nil {
		return 0, fmt.Errorf("Not ready")
	}
	return e.stdout.Read(p)
}

func (e *Epgdump) Close() error {
	if e.stdin != nil {
		e.stdin.Close()
		e.stdin = nil
	} else if e.stdout != nil {
		e.stdout.Close()
		e.stdout = nil
		if err := e.cmd.Wait(); err != nil {
			log.Print(err)
		}
	}
	return nil
}
