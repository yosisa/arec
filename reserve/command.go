package reserve

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var Recpt1Path = "/usr/local/bin/recpt1"

type Recpt1 struct {
	Path     string
	Channel  string
	Sid      string
	cmd      *exec.Cmd
	timer    *time.Timer
	cancelCh chan bool
}

func NewRecpt1(channel, sid string) *Recpt1 {
	pt1 := new(Recpt1)
	pt1.Path = Recpt1Path
	pt1.Channel = channel
	pt1.Sid = sid
	pt1.cancelCh = make(chan bool)
	return pt1
}

func (pt1 *Recpt1) Start() (io.Reader, error) {
	args := make([]string, 0)
	switch pt1.Sid {
	case "":
		args = append(args, "--b25", "--strip")
	case "epg":
		args = append(args, "--sid", pt1.Sid)
	default:
		args = append(args, "--b25", "--strip", "--sid", pt1.Sid)
	}
	args = append(args, pt1.Channel, "-", "-")
	log.Print(args)
	pt1.cmd = exec.Command(pt1.Path, args...)

	stdout, err := pt1.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := pt1.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go io.Copy(os.Stderr, stderr)

	if err := pt1.cmd.Start(); err != nil {
		return nil, err
	}
	return stdout, nil
}

func (pt1 *Recpt1) Close() error {
	if err := pt1.cmd.Process.Signal(os.Interrupt); err != nil {
		log.Print(err)
	}
	pt1.cmd.Wait()
	return nil
}

func (pt1 *Recpt1) CloseAfter(duration time.Duration) {
	log.Printf("recpt1 will be closed after %v", duration)
	pt1.timer = time.NewTimer(duration)

	go func() {
		select {
		case <-pt1.timer.C:
			pt1.Close()
		case <-pt1.cancelCh:
			log.Printf("Closing recpt1 has been canceled")
		}
	}()
}
