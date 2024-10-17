package main

import (
	"github.com/kohmebot/livebili/livebili"
	"github.com/kohmebot/plugin"
)

func NewPlugin() plugin.Plugin {
	return livebili.NewPlugin()
}
