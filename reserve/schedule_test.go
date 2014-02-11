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

func TestRemoveFromTimeline(t *testing.T) {
	tl := NewTimeline()
	p1 := &RecordInfo{Id: "1", Start: 50, End: 100}
	p2 := &RecordInfo{Id: "2", Start: 200, End: 250}
	p3 := &RecordInfo{Id: "3", Start: 150, End: 200}
	p4 := &RecordInfo{Id: "4", Start: 100, End: 150}

	tl.Set(p1)
	tl.Set(p2)
	tl.Set(p3)
	tl.Set(p4)
	assert.Equal(t, tl, &Timeline{p1, p4, p3, p2})

	assert.Nil(t, tl.Remove(p1))
	assert.Equal(t, tl, &Timeline{p4, p3, p2})
	assert.Error(t, tl.Remove(p1), "Not found")
	assert.Equal(t, tl, &Timeline{p4, p3, p2})
	assert.Nil(t, tl.Remove(p3))
	assert.Equal(t, tl, &Timeline{p4, p2})
	assert.Nil(t, tl.Remove(p2))
	assert.Equal(t, tl, &Timeline{p4})
}

func TestScheduler(t *testing.T) {
	s := NewScheduler(2, 2)
	gr1 := &RecordInfo{Id: "1", Type: "GR", Ch: "1", Start: 50, End: 100}
	gr2 := &RecordInfo{Id: "2", Type: "GR", Ch: "1", Start: 100, End: 150}
	gr3 := &RecordInfo{Id: "3", Type: "GR", Ch: "2", Start: 50, End: 150}
	gr4 := &RecordInfo{Id: "4", Type: "GR", Ch: "3", Start: 100, End: 150}

	bs1 := &RecordInfo{Id: "5", Type: "BS", Ch: "1", Start: 0, End: 100}
	bs2 := &RecordInfo{Id: "6", Type: "BS", Ch: "2", Start: 50, End: 150}
	bs3 := &RecordInfo{Id: "7", Type: "BS", Ch: "1", Start: 100, End: 200}
	bs4 := &RecordInfo{Id: "8", Type: "BS", Ch: "3", Start: 120, End: 160}

	assert.Nil(t, s.Reserve(gr1))
	assert.Nil(t, s.Reserve(gr2))
	assert.Nil(t, s.Reserve(gr3))
	assert.Error(t, s.Reserve(gr4), "Not reserved")

	assert.Nil(t, s.Reserve(bs1))
	assert.Nil(t, s.Reserve(bs2))
	assert.Nil(t, s.Reserve(bs3))
	assert.Error(t, s.Reserve(bs4), "Not reserved")

	assert.Nil(t, s.Cancel(gr2))
	assert.Error(t, s.Cancel(gr2), "Not reserved")
	assert.Nil(t, s.Reserve(gr4))

	assert.Nil(t, s.Cancel(bs2))
	assert.Error(t, s.Cancel(bs2), "Not reserved")
	assert.Nil(t, s.Reserve(bs4))
}
