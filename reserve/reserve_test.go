package reserve

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
	"testing"
)

func TestChannelEqual(t *testing.T) {
	c1 := Channel{
		Id:   "GR0_01",
		Name: "Test TV1",
	}
	c2 := c1
	assert.True(t, c1.Equal(&c2))

	c2.Name = "Test TV2"
	assert.False(t, c1.Equal(&c2))
}

func TestProgramEqual(t *testing.T) {
	p1 := Program{
		EventId:  1,
		Title:    "番組1",
		Detail:   "description here",
		Start:    1388653200,
		End:      1388654100,
		Duration: 900,
	}
	p2 := p1
	assert.True(t, p1.Equal(&p2))

	p2.Id = bson.NewObjectId()
	assert.True(t, p1.Equal(&p2))

	p2.EventId = 2
	assert.False(t, p1.Equal(&p2))
}
