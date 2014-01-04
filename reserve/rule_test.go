package reserve

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
	"testing"
)

func TestRuleMakeHash(t *testing.T) {
	r1 := Rule{Keyword: "test1"}
	r2 := Rule{Keyword: "test2"}
	assert.NotEqual(t, r1.MakeHash(), r2.MakeHash())
}

func TestRuleMakeQuery(t *testing.T) {
	id := bson.NewObjectId()
	r := Rule{Id: id, Keyword: "test"}
	assert.Equal(t, r.MakeQuery(0),
		bson.M{"reserved_by": bson.M{"$ne": id},
			"title": bson.RegEx{"test", "i"}})
	assert.Equal(t, r.MakeQuery(1388653200),
		bson.M{"reserved_by": bson.M{"$ne": id},
			"title":      bson.RegEx{"test", "i"},
			"updated_at": 1388653200})
}
