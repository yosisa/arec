package reserve

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeline(t *testing.T) {
	tl := NewTimeline()
	p1 := &RecordInfo{Start: 50, End: 100}
	p2 := &RecordInfo{Start: 200, End: 250}
	p3 := &RecordInfo{Start: 150, End: 200}
	p4 := &RecordInfo{Start: 100, End: 150}

	assert.Equal(t, tl.Set(p1), nil)
	assert.Equal(t, tl.Set(p2), nil)
	assert.Equal(t, tl.Set(p3), nil)
	assert.Error(t, tl.Set(&RecordInfo{Start: 50, End: 51}), "Duplicated schedule")
	assert.Error(t, tl.Set(&RecordInfo{Start: 99, End: 100}), "Duplicated schedule")
	assert.Error(t, tl.Set(&RecordInfo{Start: 0, End: 250}), "Duplicated schedule")
	assert.Equal(t, tl.Set(p4), nil)

	assert.Equal(t, tl.Len(), 4)
	assert.Equal(t, (*tl)[0], p1)
	assert.Equal(t, (*tl)[1], p4)
	assert.Equal(t, (*tl)[2], p3)
	assert.Equal(t, (*tl)[3], p2)

	assert.Equal(t, heap.Pop(tl), p1)
	assert.Equal(t, heap.Pop(tl), p4)
	assert.Equal(t, heap.Pop(tl), p3)
	assert.Equal(t, heap.Pop(tl), p2)
	assert.Equal(t, tl.Len(), 0)
}

func TestScheduler(t *testing.T) {
	s := NewScheduler(2, 2)
	gr1 := &RecordInfo{Ch: "GR0_1", Start: 50, End: 100}
	gr2 := &RecordInfo{Ch: "GR0_1", Start: 100, End: 150}
	gr3 := &RecordInfo{Ch: "GR0_2", Start: 50, End: 150}
	gr4 := &RecordInfo{Ch: "GR0_3", Start: 100, End: 150}

	bs1 := &RecordInfo{Ch: "BS_1", Start: 0, End: 100}
	bs2 := &RecordInfo{Ch: "BS_2", Start: 50, End: 150}
	bs3 := &RecordInfo{Ch: "BS_1", Start: 100, End: 200}
	bs4 := &RecordInfo{Ch: "BS_3", Start: 120, End: 160}

	assert.Nil(t, s.Reserve(gr1))
	assert.Nil(t, s.Reserve(gr2))
	assert.Nil(t, s.Reserve(gr3))
	assert.Error(t, s.Reserve(gr4), "Not Reserved")

	assert.Nil(t, s.Reserve(bs1))
	assert.Nil(t, s.Reserve(bs2))
	assert.Nil(t, s.Reserve(bs3))
	assert.Error(t, s.Reserve(bs4), "Not Reserved")
}
