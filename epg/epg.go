package epg

import (
	"code.google.com/p/go.text/unicode/norm"
	"encoding/json"
	"fmt"
	"github.com/yosisa/arec/reserve"
	"io"
	"strings"
	"time"
)

var flagMap map[string]string = map[string]string{
	"新": "new",
	"終": "final",
	"再": "rerun",
}

type Channel struct {
	Id        string
	Name      string
	ServiceId int `json:"service_id"`
	Programs  []Program
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

func SaveEPG(r io.Reader, ch string) error {
	channels, err := DecodeJson(r)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, channel := range channels {
		channel.Save(ch)
		for _, program := range channel.Programs {
			program.Save(now)
		}
	}
	reserve.ApplyAllRules(int(now))

	return nil
}

func (self *Channel) Save(ch string) error {
	return self.toDocument(ch).Save()
}

func (self *Channel) toDocument(ch string) *reserve.Channel {
	return &reserve.Channel{
		Id:   self.Id,
		Name: self.Name,
		Type: self.Id[:2],
		Ch:   ch,
		Sid:  self.ServiceId,
	}
}

func (self *Program) Save(now int64) error {
	return self.toDocument(now).Save()
}

func (self *Program) toDocument(now int64) *reserve.Program {
	category := make([]string, 0)
	for _, cat := range self.Categories {
		category = append(category, cat.Large.Ja, cat.Large.En, cat.Middle.Ja, cat.Middle.En)
	}

	program := &reserve.Program{
		Channel:   self.Channel,
		EventId:   fmt.Sprintf("epg:%s:%d", self.Channel, self.EventId),
		Title:     self.Title,
		Detail:    self.Detail,
		Category:  category,
		Start:     int(self.Start / 10000),
		End:       int(self.End / 10000),
		Duration:  self.Duration,
		UpdatedAt: int(now),
	}

	for key, flag := range flagMap {
		if symbol := "【" + key + "】"; strings.Contains(program.Title, symbol) {
			program.Title = strings.Replace(program.Title, symbol, "", -1)
			program.Flag = append(program.Flag, flag)
		}
	}

	return program
}
