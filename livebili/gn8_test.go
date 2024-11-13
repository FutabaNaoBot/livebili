package livebili

import (
	"github.com/kohmebot/pkg/chain"
	"github.com/stretchr/testify/assert"
	"github.com/wdvxdr1123/ZeroBot/message"
	"testing"
)

func TestDeleteAtAll(t *testing.T) {

	var msg chain.MessageChain

	msg.Split(
		message.AtAll(),
		message.Text("hello"),
		message.Text("world"),
	)

	DeleteAtAll(&msg)
	assert.Equal(t, 3, len(msg))
	assert.Equal(t, "hello", msg[0].Data["text"])

}
