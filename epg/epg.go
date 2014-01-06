package epg

import (
	"bytes"
	"code.google.com/p/go.text/unicode/norm"
	"encoding/json"
	"fmt"
	"github.com/yosisa/arec/reserve"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

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

	return nil
}

func GetEPG(recpt1 string, epgdump string, ch string) (*bytes.Reader, error) {
	log.Printf("Get EPG data: %s", ch)
	recpt1Cmd := exec.Command(recpt1, "--sid", "epg", ch, "90", "-")
	epgdumpCmd := exec.Command(epgdump, "json", "-", "-")

	recpt1Out, err := recpt1Cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	recpt1Err, err := recpt1Cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go io.Copy(os.Stdout, recpt1Err)

	epgdumpIn, err := epgdumpCmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	epgdumpOut, err := epgdumpCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	epgdumpErr, err := epgdumpCmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go io.Copy(os.Stdout, epgdumpErr)

	if err := recpt1Cmd.Start(); err != nil {
		return nil, err
	}

	if err := epgdumpCmd.Start(); err != nil {
		return nil, err
	}

	// redirect stdout of recpt1 to stdin of epgdump
	io.Copy(epgdumpIn, recpt1Out)

	if err := recpt1Cmd.Wait(); err != nil {
		log.Print(err)
	}

	// close stdin explicitly because epgdump's pipes not close automatically by exiting recpt1
	epgdumpIn.Close()
	// read all from pipe to avoid blocking by wait method
	epgdata, _ := ioutil.ReadAll(epgdumpOut)
	epg := bytes.NewReader(epgdata)

	if err := epgdumpCmd.Wait(); err != nil {
		log.Print(err)
	}

	return epg, nil
}

func (self *Channel) Save(ch string) error {
	return self.toDocument(ch).Save()
}

func (self *Channel) toDocument(ch string) *reserve.Channel {
	return &reserve.Channel{
		Id:   self.Id,
		Name: self.Name,
		Ch:   ch,
		Sid:  self.ServiceId,
	}
}

func (self *Program) Save(now int64) error {
	return self.toDocument(now).Save()
}

func (self *Program) toDocument(now int64) *reserve.Program {
	return &reserve.Program{
		Channel:   self.Channel,
		EventId:   fmt.Sprintf("epg:%s:%d", self.Channel, self.EventId),
		Title:     self.Title,
		Detail:    self.Detail,
		Start:     int(self.Start / 10000),
		End:       int(self.End / 10000),
		Duration:  self.Duration,
		UpdatedAt: int(now),
	}
}
