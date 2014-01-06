package epg

import (
	"github.com/stretchr/testify/assert"
	"github.com/yosisa/arec/reserve"
	"os"
	"testing"
	"time"
)

func TestDeocdeJson(t *testing.T) {
	f, err := os.Open("testdata/gr99.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	channels, err := DecodeJson(f)
	ch := channels[0]
	assert.Nil(t, err)
	assert.Equal(t, ch.Id, "GR0_9")
	assert.Equal(t, ch.Name, "Test TV1")
	assert.Equal(t, ch.Programs[0], Program{100, "GR0_9", "番組1", "description here",
		nil, 1.3886532e+13, 1.3886541e+13, 900, []category{
			category{categoryItem{"アニメ/特撮", "anime"},
				categoryItem{"国内アニメ", "Japanese animation"}}}})
}

func TestChannelToDocument(t *testing.T) {
	f, err := os.Open("testdata/gr99.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	channels, err := DecodeJson(f)
	channel := channels[0]
	ch := channel.toDocument("1")
	assert.Equal(t, ch, &reserve.Channel{
		Id:   "GR0_9",
		Name: "Test TV1",
		Ch:   "1",
		Sid:  9,
	})
}

func TestProgramToDocument(t *testing.T) {
	f, err := os.Open("testdata/gr99.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	channels, err := DecodeJson(f)
	program := channels[0].Programs[0]
	now := time.Now().Unix()
	pg := program.toDocument(now)
	assert.Equal(t, pg, &reserve.Program{
		Channel:   "GR0_9",
		EventId:   "epg:GR0_9:100",
		Title:     "番組1",
		Detail:    "description here",
		Category:  []string{"アニメ/特撮", "anime", "国内アニメ", "Japanese animation"},
		Start:     1388653200,
		End:       1388654100,
		Duration:  900,
		UpdatedAt: int(now),
	})
}
