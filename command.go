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

	EPG struct {
		Channel int      `goptions:"--ch, description='Get specified channel only'"`
		File    *os.File `goptions:"--file, rdonly, description='Feed from json file'"`
	} `goptions:"epg"`
	Scheduler struct{} `goptions:"scheduler"`
	Rule      struct {
		Keyword string `goptions:"--keyword, obligatory, description='regex for title'"`
	} `goptions:"rule"`
}

type SubCommand func(options *CmdOptions, config *Config)

var commands map[string]SubCommand

func SchedulerCommand(options *CmdOptions, config *Config) {
	scheduler := reserve.NewScheduler()
	scheduler.RunForever()
}

func RuleCommand(options *CmdOptions, config *Config) {
	rule := reserve.Rule{Keyword: options.Rule.Keyword}
	if err := rule.Save(); err != nil {
		log.Fatal(err)
	}
	rule.Apply(0)
}

func EPGCommand(options *CmdOptions, config *Config) {
	if options.EPG.File != nil {
		defer options.EPG.File.Close()
		if err := epg.SaveEPG(options.EPG.File); err != nil {
			log.Fatal(err)
		}
	}

	for _, channel := range config.Channels["GR"] {
		r, err := epg.GetEPG(config.Recpt1, config.Epgdump, channel.Ch)
		if err != nil {
			log.Print(err)
		} else {
			if err := epg.SaveEPG(r); err != nil {
				log.Fatal(err)
			}
		}
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
