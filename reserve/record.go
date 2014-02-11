package reserve

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RecorderItem interface {
	Info() *RecordInfo
	io.WriteCloser
}

type Engine struct {
	Scheduler *Scheduler
}

func NewEngine(gr, bs int) *Engine {
	e := new(Engine)
	e.Scheduler = NewScheduler(gr, bs)
	return e
}

func (e *Engine) Reserve(item RecorderItem) error {
	info := item.Info()
	if err := e.Scheduler.Reserve(info); err != nil {
		return err
	}
	log.Printf("%+v scheduled to record", *info)

	now := time.Now().Unix()
	wait := info.Start - int(now)
	if wait > 0 {
		info.timer = time.NewTimer(time.Duration(wait) * time.Second)
		log.Printf("Recording for %s scheduled after %d seconds", info.Id, wait)
		go func() {
			select {
			case <-info.timer.C:
				e.Record(item)
			case <-info.cancelCh:
			}
		}()
	} else if rest := wait + (info.End - info.Start); rest > 0 {
		log.Printf("Recording for %s is starting immediately", info.Id)
		go e.Record(item)
	} else {
		log.Printf("Program %s is already finished", info.Id)
	}

	return nil
}

func (e *Engine) Record(item RecorderItem) {
	info := item.Info()
	recpt1 := NewRecpt1(info.Ch, info.Sid)
	duration := time.Unix(int64(info.End), 0).Sub(time.Now())
	recpt1.CloseAfter(duration)
	io.Copy(item, recpt1)
	item.Close()
}

func (e *Engine) RunForever() {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		switch <-signalCh {
		case syscall.SIGHUP:
			log.Printf("Rescheduling")
		default:
			os.Exit(0)
		}
	}
}
