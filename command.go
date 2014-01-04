package main

import (
	"flag"
	"github.com/rakyll/command"
	"github.com/yosisa/arec/epg"
	"github.com/yosisa/arec/reserve"
	"log"
	"os"
)

var configFile *string = flag.String("config", "./arec.json", "path to config file")

type EPGCommand struct {
}

func (cmd *EPGCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (cmd *EPGCommand) Run(args []string) {
	if err := epg.SaveEPG(os.Stdin); err != nil {
		log.Fatal(err)
	}

	rule := reserve.Rule{Keyword: "news"}
	if err := rule.Save(); err != nil {
		log.Print(err)
	}
	reserve.ApplyAllRules(0)
}

type SchedulerCommand struct {
}

func (cmd *SchedulerCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (cmd *SchedulerCommand) Run(args []string) {
	scheduler := reserve.NewScheduler()
	scheduler.RunForever()
}

func init() {
	command.On("epg", "get and update epg data", &EPGCommand{})
	command.On("scheduler", "record programs", &SchedulerCommand{})
}
