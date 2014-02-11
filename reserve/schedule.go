package reserve

import (
	"container/heap"
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RecordInfo struct {
	Type     string
	Id       string
	Ch       string
	Sid      string
	Start    int
	End      int
	timer    *time.Timer
	cancelCh chan bool
}

type Timeline []*RecordInfo

func NewTimeline() *Timeline {
	t := new(Timeline)
	heap.Init(t)
	return t
}

func (t Timeline) Len() int {
	return len(t)
}

func (t Timeline) Less(i, j int) bool {
	return t[i].Start < t[j].Start
}

func (t Timeline) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *Timeline) Push(x interface{}) {
	*t = append(*t, x.(*RecordInfo))
}

func (t *Timeline) Pop() interface{} {
	old := *t
	n := len(old)
	x := old[n-1]
	*t = old[:n-1]
	return x
}

func (t *Timeline) Set(ri *RecordInfo) error {
	for _, item := range *t {
		if item.Start >= ri.End {
			break
		}
		if ri.Start-item.End < 0 {
			return fmt.Errorf("Duplicated schedule")
		}
	}
	heap.Push(t, ri)
	return nil
}

func (t *Timeline) Remove(ri *RecordInfo) error {
	for i, item := range *t {
		if item.Id == ri.Id {
			old := *t
			*t = append(old[:i], old[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Not found")
}

type Scheduler struct {
	GR []*Timeline
	BS []*Timeline
}

func NewScheduler(gr int, bs int) *Scheduler {
	tuner := new(Scheduler)
	tuner.GR = make([]*Timeline, gr)
	tuner.BS = make([]*Timeline, bs)
	for i, _ := range tuner.GR {
		tuner.GR[i] = NewTimeline()
	}
	for i, _ := range tuner.BS {
		tuner.BS[i] = NewTimeline()
	}

	return tuner
}

func (t *Scheduler) Reserve(ri *RecordInfo) error {
	for _, tl := range t.getTimelines(ri.Type) {
		if err := tl.Set(ri); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Not reserved")
}

func (t *Scheduler) Cancel(ri *RecordInfo) error {
	for _, tl := range t.getTimelines(ri.Type) {
		if err := tl.Remove(ri); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Not reserved")
}

func (t *Scheduler) getTimelines(type_ string) []*Timeline {
	if type_ == "GR" {
		return t.GR
	} else {
		return t.BS
	}
}

type Schedule struct {
	EventId   string
	StartTime int
	Duration  int
	timer     *time.Timer
	cancelCh  chan bool
}

func NewSchedule() *Schedule {
	schedule := new(Schedule)
	schedule.cancelCh = make(chan bool)
	return schedule
}

func (self *Schedule) Start() {
	now := time.Now().Unix()
	wait := self.StartTime - int(now)
	if wait > 0 {
		self.timer = time.NewTimer(time.Duration(wait) * time.Second)
		log.Printf("Recording for %s scheduled after %d seconds", self.EventId, wait)
		go func() {
			select {
			case <-self.timer.C:
				self.Record()
			case <-self.cancelCh:
			}
		}()
	} else if rest := wait + self.Duration; rest > 0 {
		self.Duration = rest
		log.Printf("Recording for %s is starting immediately", self.EventId)
		go self.Record()
	} else {
		log.Printf("Program %s is already finished", self.EventId)
	}
}

func (self *Schedule) Cancel() {
	self.cancelCh <- true
}

func (self *Schedule) Record() {
	log.Printf("Start recording: %s, duration: %d", self.EventId, self.Duration)
	// fake recording
	time.Sleep(time.Duration(self.Duration) * time.Second)
	log.Printf("Finish recording: %s", self.EventId)
}

type Recorder struct {
	activeItems map[string]*Schedule
}

func NewRecorder() *Recorder {
	scheduler := new(Recorder)
	scheduler.activeItems = make(map[string]*Schedule)
	return scheduler
}

func (self *Recorder) Refresh() {
	var program Program
	collection := getCollection("program")
	query := bson.M{
		"reserved_by": bson.M{"$not": bson.M{"$size": 0}},
		"end":         bson.M{"$gt": int(time.Now().Unix())},
	}
	iter := collection.Find(query).Sort("start").Iter()
	defer iter.Close()
	for iter.Next(&program) {
		if _, ok := self.activeItems[program.EventId]; !ok {
			s := NewSchedule()
			s.EventId = program.EventId
			s.StartTime = program.Start
			s.Duration = program.Duration
			s.Start()
			self.activeItems[s.EventId] = s
		}
	}
}

func (self *Recorder) Shutdown(force bool) {
	for _, schedule := range self.activeItems {
		schedule.Cancel()
	}
}

func (self *Recorder) RunForever() {
	self.Refresh()
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		switch <-signalCh {
		case syscall.SIGHUP:
			log.Printf("Rescheduling")
			self.Refresh()
		default:
			os.Exit(0)
		}
	}
}
