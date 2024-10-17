package main

import (
	"github.com/futabanaobot/livebili/livebili"
	"github.com/futabanaobot/plugin"
)

func NewPlugin() plugin.Plugin {
	return livebili.NewPlugin()
}
