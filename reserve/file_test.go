package reserve

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileRecord(t *testing.T) {
	channel := Channel{
		Type: "GR",
		Ch:   "1",
		Sid:  99,
	}
	program := Program{
		Title: "program 1",
		Start: 1388590200,
		End:   1388592000,
	}
	program.Hash = program.MakeHash()
	f := NewFileRecord(channel, program)

	assert.Equal(t, f.channel, &channel)
	assert.Equal(t, f.program, &program)

	assert.Equal(t, f.Info(), &RecordInfo{
		Id:    "2e9303774d617e707cb068aa715a0050",
		Type:  "GR",
		Ch:    "1",
		Sid:   "99",
		Start: 1388590200,
		End:   1388592000,
	})

	assert.Equal(t, f.getFilename(), "20140102T0030_GR1_1800_program 1.ts")

	ChangeFilename("{{.Type}}{{.Ch}}_{{.Start}}-{{.End}}.ts")
	DateFormat = "0102_1504"
	SaveDirectory = "/var"
	assert.Equal(t, f.getFilename(), "/var/GR1_0102_0030-0102_0100.ts")
	SaveDirectory = "/var/"
	assert.Equal(t, f.getFilename(), "/var/GR1_0102_0030-0102_0100.ts")
}
