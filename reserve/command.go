package reserve

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

type Recpt1 struct {
	Path     string
	Channel  string
	Sid      string
	cmd      *exec.Cmd
	timer    *time.Timer
	cancelCh chan bool
}

func NewRecpt1(path, channel, sid string) *Recpt1 {
	pt1 := new(Recpt1)
	pt1.Path = path
	pt1.Channel = channel
	pt1.Sid = sid
	pt1.cancelCh = make(chan bool)
	return pt1
}

func (pt1 *Recpt1) Start(w io.WriteCloser) error {
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
		return err
	}

	stderr, err := pt1.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, stderr)

	if err := pt1.cmd.Start(); err != nil {
		return err
	}

	go func() {
		io.Copy(w, stdout)
		w.Close()
	}()
	return nil
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
