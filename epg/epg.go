package epg

import (
	"code.google.com/p/go.text/unicode/norm"
	"encoding/json"
	"io"
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
