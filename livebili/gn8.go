package livebili

import (
	"github.com/kohmebot/gn8/gn8sdk"
	"github.com/kohmebot/pkg/chain"
	"strings"
	"time"
)

type gn8 struct {
	i *gn8sdk.Invoker
}

func (g *gn8) IsNowDND() bool {
	return g.IfNotNil(func() bool {
		return g.i.IsDND(time.Now())
	})
}
func (g *gn8) IsDND(ti time.Time) bool {
	return g.IfNotNil(func() bool {
		return g.i.IsDND(ti)
	})
}

func (g *gn8) IfNotNil(fn func() bool) bool {
	if g.i == nil {
		return false
	}
	return fn()
}

func DeleteAtAll(c *chain.MessageChain) {
	targetIdx := -1
	for idx, segment := range *c {
		if segment.Type == "at" && segment.Data["qq"] == "all" {
			targetIdx = idx
			break
		}
	}
	if targetIdx == -1 {
		return
	}

	*c = append((*c)[:targetIdx], (*c)[targetIdx+1:]...)

	nextIdx := targetIdx

	if nextIdx >= len(*c) {
		return
	}

	seg := (*c)[nextIdx]
	if seg.Type != "text" {
		return
	}
	if len(strings.TrimSpace(seg.Data["text"])) > 0 {
		return
	}

	*c = append((*c)[:nextIdx], (*c)[nextIdx+1:]...)

}
