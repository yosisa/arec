package reserve

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleMakeHash(t *testing.T) {
	r1 := Rule{Keyword: "test1"}
	r2 := Rule{Keyword: "test2"}
	assert.NotEqual(t, r1.MakeHash(), r2.MakeHash())
}
