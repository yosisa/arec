package epg

import (
	"github.com/stretchr/testify/assert"
	"github.com/yosisa/arec/reserve"
	"os"
	"testing"
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
	assert.Equal(t, ch.Id, "GR0_01")
	assert.Equal(t, ch.Name, "Test TV1")
	assert.Equal(t, ch.Programs[0], Program{100, "GR0_01", "番組1", "description here",
		nil, 1.3886532e+13, 1.3886541e+13, 900, []category{
			category{categoryItem{"アニメ/特撮", "anime"},
				categoryItem{"国内アニメ", "Japanese animation"}}}})
}

func TestProgramToDocument(t *testing.T) {
	f, err := os.Open("testdata/gr99.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	channels, err := DecodeJson(f)
	program := channels[0].Programs[0]
	pg := program.toDocument()
	assert.Equal(t, pg, &reserve.Program{
		EventId:  100,
		Title:    "番組1",
		Detail:   "description here",
		Start:    1388653200,
		End:      1388654100,
		Duration: 900,
	})
}
