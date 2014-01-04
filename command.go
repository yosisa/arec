package main

import (
	"github.com/voxelbrain/goptions"
	"github.com/yosisa/arec/epg"
	"github.com/yosisa/arec/reserve"
	"log"
	"os"
)

type CmdOptions struct {
	Config string        `goptions:"-c, --config, description='Path to config file'"`
	Help   goptions.Help `goptions:"-h, --help, description='Show this help'"`
	goptions.Verbs

	EPG       struct{} `goptions:"epg"`
	Scheduler struct{} `goptions:"scheduler"`
	Rule      struct {
		Keyword string `goptions:"--keyword, obligatory, description='regex for title'"`
	} `goptions:"rule"`
}

type SubCommand func(options *CmdOptions)

var commands map[string]SubCommand

func SchedulerCommand(options *CmdOptions) {
	scheduler := reserve.NewScheduler()
	scheduler.RunForever()
}

func RuleCommand(options *CmdOptions) {
	rule := reserve.Rule{Keyword: options.Rule.Keyword}
	if err := rule.Save(); err != nil {
		log.Fatal(err)
	}
	rule.Apply(0)
}

func EPGCommand(options *CmdOptions) {
	if err := epg.SaveEPG(os.Stdin); err != nil {
		log.Fatal(err)
	}

	rule := reserve.Rule{Keyword: "news"}
	if err := rule.Save(); err != nil {
		log.Print(err)
	}
	reserve.ApplyAllRules(0)
}

func init() {
	commands = map[string]SubCommand{
		"scheduler": SchedulerCommand,
		"rule":      RuleCommand,
		"epg":       EPGCommand,
	}
}
