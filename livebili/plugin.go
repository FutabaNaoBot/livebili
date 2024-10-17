package livebili

import (
	"github.com/futabanaobot/plugin"
	"github.com/futabanaobot/plugin/pkg/command"
	"github.com/futabanaobot/plugin/pkg/version"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type biliPlugin struct {
	e    *zero.Engine
	env  plugin.Env
	conf Config
}

func NewPlugin() plugin.Plugin {
	return &biliPlugin{}
}

func (b *biliPlugin) Init(engine *zero.Engine, env plugin.Env) error {
	b.e = engine
	b.env = env

	return b.init()
}

func (b *biliPlugin) Name() string {
	return "BiliBili-Live"
}

func (b *biliPlugin) Description() string {
	return "推送开播信息"
}

func (b *biliPlugin) Commands() command.Commands {
	return command.NewCommands()
}

func (b *biliPlugin) Version() version.Version {
	return version.NewVersion(0, 0, 1)
}
