package epg

import (
	"code.google.com/p/go.text/unicode/norm"
	"encoding/json"
	"github.com/yosisa/arec/reserve"
	"io"
	"time"
)

type Channel struct {
	Id       string
	Name     string
	Programs []Program
}

type Program struct {
	EventId    int `json:"event_id"`
	Channel    string
	Title      string
	Detail     string
	ExtDetails []extDetail `json:"extdetail"`
	Start      float64
	End        float64
	Duration   int
	Categories []category `json:"category"`
}

type extDetail struct {
	Key   string `json:"item_description"`
	Value string `json:"item"`
}

type category struct {
	Large  categoryItem
	Middle categoryItem
}

type categoryItem struct {
	Ja string `json:"ja_JP"`
	En string
}

func DecodeJson(r io.Reader) ([]Channel, error) {
	rr := norm.NFKC.Reader(r)
	dec := json.NewDecoder(rr)
	var channels []Channel
	return channels, dec.Decode(&channels)
}

func SaveEPG(r io.Reader) error {
	channels, err := DecodeJson(r)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, channel := range channels {
		channel.Save()
		for _, program := range channel.Programs {
			program.Save(now)
		}
	}

	return nil
}

func (self *Channel) Save() error {
	ch := reserve.Channel{self.Id, self.Name}
	return ch.Save()
}

func (self *Program) Save(now int64) error {
	return self.toDocument(now).Save()
}

func (self *Program) toDocument(now int64) *reserve.Program {
	return &reserve.Program{
		EventId:   self.EventId,
		Title:     self.Title,
		Detail:    self.Detail,
		Start:     int(self.Start / 10000),
		End:       int(self.End / 10000),
		Duration:  self.Duration,
		UpdatedAt: int(now),
	}
}
