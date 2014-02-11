package epg

import (
	"github.com/yosisa/arec/command"
	"github.com/yosisa/arec/reserve"
	"log"
	"time"
)

const (
	GR_REC_TIME = 90 * time.Second
	BS_REC_TIME = 300 * time.Second
)

type EpgRecord struct {
	recordInfo *reserve.RecordInfo
	command.Epgdump
}

func NewEpgRecord(type_, ch string, start time.Time) *EpgRecord {
	var duration time.Duration
	if type_ == "GR" {
		duration = GR_REC_TIME
	} else {
		duration = BS_REC_TIME
	}

	epg := new(EpgRecord)
	epg.recordInfo = &reserve.RecordInfo{
		Type:  type_,
		Ch:    ch,
		Sid:   "epg",
		Start: int(start.Unix()),
		End:   int(start.Add(duration).Unix()),
	}

	return epg
}

func (e *EpgRecord) Info() *reserve.RecordInfo {
	return e.recordInfo
}

func (e *EpgRecord) Close() error {
	if err := e.Epgdump.Close(); err != nil {
		log.Println(err)
		return err
	}

	SaveEPG(e, e.recordInfo.Ch)
	return e.Epgdump.Close()
}

func Reserve(engine *reserve.Engine, type_, ch string) {
	var duration time.Duration
	if type_ == "GR" {
		duration = GR_REC_TIME
	} else {
		duration = BS_REC_TIME
	}
	start := time.Now()

	for {
		epgRecord := NewEpgRecord(type_, ch, start)
		if err := engine.Reserve(epgRecord); err == nil {
			return
		}
		start = start.Add(duration)
	}
}
